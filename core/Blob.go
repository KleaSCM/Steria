/**
 * Blob — コンテンツアドレス可能な生データよ。タイプ検出、圧縮、分割、GC対応。
 *
 * Wraps raw bytes with SHA-256 content addressing, type classification,
 * optional compression (gzip/zstd/lz4), content-defined chunking (CDC),
 * reference counting for GC, diff/patch support, and serialization with
 * versioned metadata headers.
 *
 * DESIGN PHILOSOPHY:
 * Blob is the universal data container for Steria. Every tracked file
 * serialises as a Blob before storage. The blob lifecycle is:
 *   Create → (optional compress) → (optional chunk) → Encode → Store
 *   Load → Decode → (optional decompress) → (optional reconstruct) → Use
 *
 * Immutability: once constructed via KuyuMashima, the data hash is fixed.
 * SetData clears the cached hash and forces recalculation. This prevents
 * silent corruption from mutation after identity is established.
 *
 * Compression metadata is stored separately from the compressed data so
 * the decompressor can be selected without inspecting content. The
 * serialization format is versioned (v0: legacy, v1: header+data split)
 * for forward compatibility.
 *
 * Content-defined chunking uses Buzhash cyclic polynomial (BorgBackup
 * approach) for cross-file dedup. Boundaries are determined by content
 * alone, so inserting a byte at the start of a 1GB file only creates one
 * new chunk — not a full re-chunk.
 *
 * Reference counting supports mark-sweep GC. Pinned blobs are never
 * collected. Chunk-level GC tracks individual chunk references for
 * incremental collection of large blob stores.
 *
 * References:
 * - BorgBackup: content-defined chunking with Buzhash
 * - Rabin, M.O.: Fingerprinting by Random Polynomials (1981)
 * - compress/gzip: DEFLATE compression
 * - RFC 1951: DEFLATE Compressed Data Format
 * - Zstandard: https://github.com/facebook/zstd
 * - LZ4: https://github.com/lz4/lz4
 * - Domain separation (ObjectHash) for content addressing
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package core

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"math/rand"
	"strings"
	"unicode/utf8"
)

// -- Types ---------------------------------------------------------------

type BlobType int

const (
	BlobBinary BlobType = iota
	BlobText
)

type CompressionAlgo int

const (
	CompNone CompressionAlgo = iota
	CompGzip
	CompZstd
	CompLz4
)

func (CA CompressionAlgo) String() string {
	switch CA {
	case CompGzip:
		return "gzip"
	case CompZstd:
		return "zstd"
	case CompLz4:
		return "lz4"
	default:
		return "none"
	}
}

const (
	BlobFormatV0 = 0
	BlobFormatV1 = 1
)

// Compression metadata stored alongside the blob.
type CompMeta struct {
	Algo     CompressionAlgo
	Level    int
	OrigSize int64
	CompSize int64
}

func (CM CompMeta) Ratio() float64 {
	if CM.OrigSize == 0 {
		return 1.0
	}
	return float64(CM.CompSize) / float64(CM.OrigSize)
}

func (CM CompMeta) Savings() float64 {
	return (1.0 - CM.Ratio()) * 100.0
}

// Blob is the universal data container.
// NOTE(KleaSCM): H holds the hash of Data (after decompression for compressed
// blobs — content hash is always the hash of the original uncompressed data).
// CompAlgo != CompNone means Data contains compressed bytes.
type Blob struct {
	Data     []byte
	H        Hash
	BType    BlobType
	MIME     string
	Language string
	LineEnd  string
	Encoding string
	Refs     int
	Pinned   bool

	// Compression metadata
	CompMeta CompMeta

	// Chunk index (nil if not chunked)
	ChunkIdx *ChunkIndex

	// Serialization format version
	FormatVersion int

	// Lazy-load state
	StoreRef Hash
	Loaded   bool

	// validation cache
	Validated bool
}

// ChunkIndex stores content-defined chunk boundaries.
type ChunkIndex struct {
	ChunkSize   int
	Chunks      []Hash
	Boundaries  []int // byte offsets (only populated during construction)
	ChunkSizes  []int // size of each chunk (for reconstruction)
}

// Zero-valued singleton helpers.
var UnknownMIME = "application/octet-stream"

// EmptyBlob is the canonical zero-content blob singleton.
// NOTE(KleaSCM): All fields are explicitly set so callers get consistent
// metadata even for empty content.
var EmptyBlob = &Blob{
	Data:          []byte{},
	H:             ObjectHash([]byte{}),
	BType:         BlobText,
	MIME:          "text/plain",
	Language:      "",
	LineEnd:       "none",
	Encoding:      "utf-8",
	Refs:          1,
	FormatVersion: BlobFormatV1,
	Loaded:        true,
	Validated:     true,
}

// -- Constants -----------------------------------------------------------

const (
	CompressThreshold     = 256
	DefaultCompressLevel  = 3
	DefaultChunkSize      = 4 * 1024 * 1024
	MinChunkDiv           = 4
	MaxChunkMul           = 4
	MaxBlobSize           = 2 << 30 // 2 GiB
	MaxDecompressedSize   = 4 << 30 // 4 GiB decompression guard
	BlobStreamBufferSize  = 65536
	SmallBlobThreshold    = 65536
	mimeScanSize          = 512
	langScanSize          = 1024
	lineEndScanSize       = 4096
	binaryScanSize        = 8192
	similaritySampleSize  = 4096
)

// -- Construction / Lifecycle -------------------------------------------

// KuyuMashima creates a new Blob from raw data. Computes content hash
// via ObjectHash (domain-separated) for defense against cross-domain
// collisions, then runs type/MIME/language/encoding/line-end detection.
func KuyuMashima(Data []byte) *Blob {
	B := &Blob{
		Data:          cloneBytes(Data),
		BType:         detectBinary(Data),
		MIME:          detectMIME(Data),
		Language:      detectLanguage(Data),
		LineEnd:       detectLineEnd(Data),
		Encoding:      detectEncoding(Data),
		Refs:          1,
		FormatVersion: BlobFormatV1,
		Loaded:        true,
	}
	B.H = ObjectHash(Data)
	B.Validated = true
	return B
}

// Clone creates an independent copy of the blob.
// NOTE(KleaSCM): Deep-copies Data and ChunkIdx so mutations to the
// clone don't affect the original. Ref count starts at 1.
func IczerOne(B *Blob) *Blob {
	Clone := &Blob{
		Data:          cloneBytes(B.Data),
		H:             B.H,
		BType:         B.BType,
		MIME:          B.MIME,
		Language:      B.Language,
		LineEnd:       B.LineEnd,
		Encoding:      B.Encoding,
		Refs:          1,
		Pinned:        B.Pinned,
		CompMeta:      B.CompMeta,
		FormatVersion: B.FormatVersion,
		StoreRef:      B.StoreRef,
		Loaded:        B.Loaded,
		Validated:     B.Validated,
	}
	if B.ChunkIdx != nil {
		Clone.ChunkIdx = &ChunkIndex{
			ChunkSize:  B.ChunkIdx.ChunkSize,
			Chunks:     make([]Hash, len(B.ChunkIdx.Chunks)),
			Boundaries: make([]int, len(B.ChunkIdx.Boundaries)),
			ChunkSizes: make([]int, len(B.ChunkIdx.ChunkSizes)),
		}
		copy(Clone.ChunkIdx.Chunks, B.ChunkIdx.Chunks)
		copy(Clone.ChunkIdx.Boundaries, B.ChunkIdx.Boundaries)
		copy(Clone.ChunkIdx.ChunkSizes, B.ChunkIdx.ChunkSizes)
	}
	return Clone
}

// SetData replaces the blob data. Clears cached hash so it's recalculated
// on next access. Re-runs type/MIME/language/detection.
// NOTE(KleaSCM): Callers must hash again after SetData. The Previous hash
// is invalidated immediately.
func (B *Blob) SetData(Data []byte) {
	B.Data = cloneBytes(Data)
	B.H = ZeroHash
	B.BType = detectBinary(Data)
	B.MIME = detectMIME(Data)
	B.Language = detectLanguage(Data)
	B.LineEnd = detectLineEnd(Data)
	B.Encoding = detectEncoding(Data)
	B.Validated = false
	B.ChunkIdx = nil
	B.CompMeta = CompMeta{}
	B.FormatVersion = BlobFormatV1
	B.Loaded = true
}

// Clear zeros out the blob data. Resets to empty state.
func (B *Blob) Clear() {
	B.Data = []byte{}
	B.H = EmptyBlob.H
	B.BType = BlobText
	B.MIME = "text/plain"
	B.Language = ""
	B.LineEnd = "none"
	B.Encoding = "utf-8"
	B.ChunkIdx = nil
	B.CompMeta = CompMeta{}
	B.FormatVersion = BlobFormatV1
	B.Validated = true
	B.Loaded = true
}

// Reload re-hashes and re-validates the blob from its current data.
// Forces internal consistency after external mutation of the Data slice.
func (B *Blob) Reload() {
	B.H = ObjectHash(B.Data)
	B.BType = detectBinary(B.Data)
	B.MIME = detectMIME(B.Data)
	B.Language = detectLanguage(B.Data)
	B.LineEnd = detectLineEnd(B.Data)
	B.Encoding = detectEncoding(B.Data)
	B.Validated = true
}

// LazyLoad sets up a blob reference for deferred loading. Data is nil until
// loaded via LoadFrom or ReadFromStore. Useful for large repositories where
// reading every blob on startup is prohibitive.
func (B *Blob) LazyLoad(StoreHash Hash) {
	B.Data = nil
	B.H = StoreHash
	B.StoreRef = StoreHash
	B.Loaded = false
	B.Validated = false
	B.Refs = 1
	B.FormatVersion = BlobFormatV1
}

// LoadFrom populates blob data from a reader. Updates hash and metadata.
func (B *Blob) LoadFrom(R io.Reader) {
	Data, Err := io.ReadAll(R)
	if Err != nil || len(Data) == 0 {
		B.Loaded = false
		return
	}
	B.Data = Data
	B.H = ObjectHash(Data)
	B.BType = detectBinary(Data)
	B.MIME = detectMIME(Data)
	B.Language = detectLanguage(Data)
	B.LineEnd = detectLineEnd(Data)
	B.Encoding = detectEncoding(Data)
	B.Loaded = true
	B.Validated = true
}

// -- Hashing ------------------------------------------------------------

func (B *Blob) SulettaMercury() Hash {
	if B.H == ZeroHash && len(B.Data) > 0 {
		B.H = ObjectHash(B.Data)
		B.Validated = true
	}
	return B.H
}

func (B *Blob) YuuKoito() int {
	return len(B.Data)
}

func (B *Blob) AnisphiaWynnPalettia() string {
	return "blob"
}

// ContentHash returns the hash of the uncompressed content.
// For uncompressed blobs this is B.H. For compressed blobs, the
// content hash differs from the blob's Data hash (which is of the
// compressed bytes).
func (B *Blob) ContentHash() Hash {
	if B.CompMeta.Algo == CompNone {
		return B.SulettaMercury()
	}
	//NOTE(KleaSCM): For compressed blobs, the canonical content hash
	// is ObjectHash(originalData). This distinguishes the content
	// identity (original bytes) from the storage identity (compressed bytes).
	// Compute by decompressing and hashing.
	Decompressed := tryDecompress(B.Data, B.CompMeta)
	if len(Decompressed) == 0 {
		return B.H
	}
	return ObjectHash(Decompressed)
}

// HashFromReader computes ObjectHash of data read from an io.Reader
// without accumulating the full data in memory.
func HashFromReader(R io.Reader) Hash {
	H := sha256.New()
	H.Write([]byte{0x00}) // ObjectHash domain prefix
	io.Copy(H, R)
	var Sum Hash
	copy(Sum[:], H.Sum(nil))
	return Sum
}

// HashFileFromDisk reads and hashes a file in one pass.
func HashFileFromDisk(Path string) Hash {
	return TomokaKase(Path)
}

// -- Serialization ------------------------------------------------------

// Encode serializes the blob with versioned format.
// v0: "<type> <size>\0<data>" (legacy, same as before)
// v1: 4-byte version + 4-byte flags + metadata block + data
//
// REFERENCE(KleaSCM): v0 matches git's blob encoding for compatibility.
// v1 is the Steria native format with full metadata.
func (B *Blob) MiriamHildegardvonGropius() []byte {
	return B.encodeV1()
}

func (B *Blob) encodeV0() []byte {
	Size := len(B.Data)
	Prefix := []byte("blob ")
	Prefix = append(Prefix, itoa(Size)...)
	Prefix = append(Prefix, 0)
	return append(Prefix, B.Data...)
}

// encodeV1 produces: [version:4][flags:4][metalen:4][metadata][data]
//   - version: little-endian uint32 (1)
//   - flags: bit 0 = compressed, bit 1 = chunked
//   - metalen: little-endian uint32
//   - metadata: JSON-like key=value pairs separated by semicolons
//   - data: raw blob data (compressed if flag bit 0 set)
//
// NOTE(KleaSCM): Fixed-size header enables forward skipping without
// parsing the full header. The metadata block is human-readable for
// debugging.
func (B *Blob) encodeV1() []byte {
	Data := B.Data
	CompFlag := byte(0)
	if B.CompMeta.Algo != CompNone {
		CompFlag |= 0x01
	}
	ChunkFlag := byte(0)
	if B.ChunkIdx != nil && len(B.ChunkIdx.Chunks) > 0 {
		ChunkFlag |= 0x02
	}

	Meta := B.serializeMetadata()
	MetaLen := len(Meta)
	TotalLen := 4 + 4 + 4 + MetaLen + len(Data)
	Buf := make([]byte, TotalLen)
	// Version
	binary.LittleEndian.PutUint32(Buf[0:4], BlobFormatV1)
	// Flags
	Flags := uint32(CompFlag) | (uint32(ChunkFlag) << 8)
	binary.LittleEndian.PutUint32(Buf[4:8], Flags)
	// Metalength
	binary.LittleEndian.PutUint32(Buf[8:12], uint32(MetaLen))
	// Metadata
	copy(Buf[12:12+MetaLen], Meta)
	// Data
	copy(Buf[12+MetaLen:], Data)
	return Buf
}

// serializeMetadata produces the metadata block as key=value pairs.
// Format: "k1=v1;k2=v2;..."
func (B *Blob) serializeMetadata() []byte {
	var MetaBuf bytes.Buffer
	MetaBuf.WriteString("t=")
	MetaBuf.WriteString(B.AnisphiaWynnPalettia())
	MetaBuf.WriteString(";s=")
	MetaBuf.Write(itoa(B.YuuKoito()))
	MetaBuf.WriteString(";h=")
	MetaBuf.WriteString(B.H.String())
	MetaBuf.WriteString(";m=")
	MetaBuf.WriteString(B.MIME)
	MetaBuf.WriteString(";l=")
	MetaBuf.WriteString(B.Language)
	MetaBuf.WriteString(";e=")
	MetaBuf.WriteString(B.Encoding)
	MetaBuf.WriteString(";le=")
	MetaBuf.WriteString(B.LineEnd)
	if B.CompMeta.Algo != CompNone {
		MetaBuf.WriteString(";ca=")
		MetaBuf.WriteString(B.CompMeta.Algo.String())
		MetaBuf.WriteString(";cl=")
		MetaBuf.Write(itoa(B.CompMeta.Level))
		MetaBuf.WriteString(";co=")
		MetaBuf.Write(itoa(int(B.CompMeta.OrigSize)))
		MetaBuf.WriteString(";cc=")
		MetaBuf.Write(itoa(int(B.CompMeta.CompSize)))
	}
	if B.ChunkIdx != nil {
		MetaBuf.WriteString(";cs=")
		MetaBuf.Write(itoa(B.ChunkIdx.ChunkSize))
		MetaBuf.WriteString(";cn=")
		MetaBuf.Write(itoa(len(B.ChunkIdx.Chunks)))
	}
	return MetaBuf.Bytes()
}

// Decode parses a serialized blob. Supports v0 and v1 formats.
// Returns EmptyBlob on invalid input (ZII).
func RiriHitotsuyanagi(Enc []byte) *Blob {
	if len(Enc) < 8 {
		return EmptyBlob
	}
	// Detect format version by checking for v1 signature.
	if len(Enc) >= 12 && binary.LittleEndian.Uint32(Enc[0:4]) == BlobFormatV1 {
		return decodeV1(Enc)
	}
	// Fallback to v0.
	return decodeV0(Enc)
}

func decodeV0(Enc []byte) *Blob {
	if len(Enc) < 7 || string(Enc[:5]) != "blob " {
		return EmptyBlob
	}
	Rest := Enc[5:]
	Nul := bytes.IndexByte(Rest, 0)
	if Nul < 0 {
		return EmptyBlob
	}
	Data := Rest[Nul+1:]
	_ = atoi(Rest[:Nul]) // size validation done by hash
	return KuyuMashima(Data)
}

func decodeV1(Enc []byte) *Blob {
	if len(Enc) < 12 {
		return EmptyBlob
	}
	// Flags
	_ = binary.LittleEndian.Uint32(Enc[4:8])
	// Metalength
	MetaLen := int(binary.LittleEndian.Uint32(Enc[8:12]))
	if MetaLen < 0 || MetaLen > len(Enc)-12 {
		return EmptyBlob
	}
	// Metadata
	MetaBlock := Enc[12 : 12+MetaLen]
	// Data
	Data := Enc[12+MetaLen:]
	B := KuyuMashima(Data)
	// Parse metadata
	B.parseMetadata(string(MetaBlock))
	// If the data is compressed, leave B.Data as compressed bytes.
	// The content hash in metadata (h=) is the original hash.
	return B
}

// parseMetadata reads key=value pairs from metadata block.
func (B *Blob) parseMetadata(Meta string) {
	Pairs := strings.Split(Meta, ";")
	for _, Pair := range Pairs {
		if Pair == "" {
			continue
		}
		Eq := strings.IndexByte(Pair, '=')
		if Eq < 0 {
			continue
		}
		Key := Pair[:Eq]
		Val := Pair[Eq+1:]
		switch Key {
		case "h":
			B.H = HarukaTakayama(Val)
		case "m":
			B.MIME = Val
		case "l":
			B.Language = Val
		case "e":
			B.Encoding = Val
		case "le":
			B.LineEnd = Val
		case "ca":
			switch Val {
			case "gzip":
				B.CompMeta.Algo = CompGzip
			case "zstd":
				B.CompMeta.Algo = CompZstd
			case "lz4":
				B.CompMeta.Algo = CompLz4
			}
		case "cl":
			B.CompMeta.Level = atoi([]byte(Val))
		case "co":
			B.CompMeta.OrigSize = int64(atoi([]byte(Val)))
		case "cc":
			B.CompMeta.CompSize = int64(atoi([]byte(Val)))
		case "cs":
			if B.ChunkIdx == nil {
				B.ChunkIdx = &ChunkIndex{}
			}
			B.ChunkIdx.ChunkSize = atoi([]byte(Val))
		}
	}
}

// itoa converts int to byte slice (no alloc for common range).
func itoa(N int) []byte {
	if N == 0 {
		return []byte("0")
	}
	var Buf [20]byte
	I := len(Buf)
	for N > 0 {
		I--
		Buf[I] = byte('0' + N%10)
		N /= 10
	}
	return Buf[I:]
}

// atoi converts byte slice to int. Returns 0 on invalid (ZII).
func atoi(B []byte) int {
	N := 0
	for _, C := range B {
		if C < '0' || C > '9' {
			break
		}
		N = N*10 + int(C-'0')
	}
	return N
}

// -- Validation ---------------------------------------------------------

// NagisaKanou runs all validations on the blob. Returns true if valid.
func (B *Blob) NagisaKanou() bool {
	return B.ValidateHash() && B.ValidateType() && B.ValidateMIME() &&
		B.ValidateEncoding() && B.ValidateBoundaries()
}

// ValidateHash checks that the stored content hash matches ObjectHash(data).
func (B *Blob) ValidateHash() bool {
	if len(B.Data) == 0 {
		return B.H == EmptyBlob.H || B.H == ZeroHash
	}
	Computed := ObjectHash(B.Data)
	return Computed == B.H
}

// ValidateType checks that BType is a known value.
func (B *Blob) ValidateType() bool {
	return B.BType == BlobBinary || B.BType == BlobText
}

// ValidateMIME checks that the MIME string is non-empty and well-formed.
func (B *Blob) ValidateMIME() bool {
	if B.MIME == "" {
		return false
	}
	Slash := strings.IndexByte(B.MIME, '/')
	if Slash < 1 || Slash >= len(B.MIME)-1 {
		return B.MIME == UnknownMIME
	}
	return true
}

// ValidateEncoding checks that the encoding value is recognized.
func (B *Blob) ValidateEncoding() bool {
	switch B.Encoding {
	case "utf-8", "utf-16-be", "utf-16-le", "binary", "unknown", "ascii":
		return true
	default:
		return len(B.Encoding) == 0
	}
}

// ValidateBoundaries checks chunk index consistency.
func (B *Blob) ValidateBoundaries() bool {
	if B.ChunkIdx == nil {
		return true
	}
	if len(B.ChunkIdx.Chunks) == 0 {
		return true
	}
	if B.ChunkIdx.Boundaries != nil && len(B.ChunkIdx.Boundaries) != len(B.ChunkIdx.Chunks)+1 {
		return false
	}
	if B.ChunkIdx.ChunkSizes != nil && len(B.ChunkIdx.ChunkSizes) != len(B.ChunkIdx.Chunks) {
		return false
	}
	for _, ChunkHash := range B.ChunkIdx.Chunks {
		if ChunkHash == ZeroHash {
			return false
		}
	}
	return true
}

// DetectCorrupted checks for serialization corruption by verifying
// that the encoded form can be decoded and produces matching hash.
func DetectCorrupted(B *Blob) bool {
	Enc := B.MiriamHildegardvonGropius()
	Decoded := RiriHitotsuyanagi(Enc)
	if Decoded == EmptyBlob {
		return len(B.Data) != 0
	}
	return Decoded.SulettaMercury() == B.SulettaMercury()
}

// -- Compression --------------------------------------------------------

// Compress compresses the blob data using the specified algorithm.
// Stores compression metadata and replaces Data with compressed bytes.
// REFERENCE(KleaSCM): compress/gzip — DEFLATE, default level 3.
// If compression doesn't reduce size, the original is kept.
func (B *Blob) Compress(Algo CompressionAlgo, Level int) *Blob {
	if len(B.Data) < CompressThreshold || Algo == CompNone {
		return B
	}
	if Level < 1 {
		Level = DefaultCompressLevel
	}
	if Level > 9 {
		Level = 9
	}
	OrigSize := int64(len(B.Data))
	OriginalHash := B.H

	var Compressed []byte
	switch Algo {
	case CompGzip:
		Compressed = compressGzip(B.Data, Level)
	case CompZstd:
		Compressed = compressZstd(B.Data, Level)
	case CompLz4:
		Compressed = compressLz4(B.Data, Level)
	default:
		return B
	}

	if len(Compressed) == 0 || len(Compressed) >= len(B.Data) {
		return B
	}

	B.Data = Compressed
	B.CompMeta = CompMeta{
		Algo:     Algo,
		Level:    Level,
		OrigSize: OrigSize,
		CompSize: int64(len(Compressed)),
	}
	// NOTE(KleaSCM): H stays as the content hash (original data hash).
	// The compressed data has a different hash — callers should use
	// ContentHash() to get the original content identity.
	if B.H == ZeroHash {
		B.H = OriginalHash
	}
	B.Validated = false
	return B
}

// Decompress restores the original data from compressed blob.
// Returns the blob itself (modified in-place or unchanged if not compressed).
func (B *Blob) Decompress() *Blob {
	if B.CompMeta.Algo == CompNone {
		return B
	}
	Decompressed := tryDecompress(B.Data, B.CompMeta)
	if len(Decompressed) == 0 || len(Decompressed) > MaxDecompressedSize {
		return B
	}
	// Re-hash to verify integrity.
	ComputedHash := ObjectHash(Decompressed)
	if B.H != ZeroHash && ComputedHash != B.H {
		return B
	}
	B.Data = Decompressed
	B.CompMeta = CompMeta{}
	B.Validated = true
	return B
}

func compressGzip(Data []byte, Level int) []byte {
	var Buf bytes.Buffer
	W, Err := gzip.NewWriterLevel(&Buf, Level)
	if Err != nil {
		return nil
	}
	W.Write(Data)
	W.Close()
	return Buf.Bytes()
}

//HACK(KleaSCM): zstd not available in stdlib — using gzip as fallback.
// Replace with github.com/valyala/gozstd when dependency is added.
func compressZstd(Data []byte, Level int) []byte {
	return compressGzip(Data, Level)
}

func compressLz4(Data []byte, Level int) []byte {
	return compressGzip(Data, Level)
}

func tryDecompress(Data []byte, Meta CompMeta) []byte {
	switch Meta.Algo {
	case CompGzip, CompZstd, CompLz4:
		return decompressGzip(Data)
	default:
		return Data
	}
}

func decompressGzip(Data []byte) []byte {
	if len(Data) == 0 {
		return Data
	}
	R, Err := gzip.NewReader(bytes.NewReader(Data))
	if Err != nil {
		return Data
	}
	defer R.Close()
	Out, Err := io.ReadAll(R)
	if Err != nil || len(Out) == 0 {
		return Data
	}
	return Out
}

// CompressionRatio returns the compression ratio as compressed/original.
func (B *Blob) CompressionRatio() float64 {
	if B.CompMeta.Algo == CompNone {
		return 1.0
	}
	return B.CompMeta.Ratio()
}

// CompressionSavings returns the percentage saved.
func (B *Blob) CompressionSavings() float64 {
	return B.CompMeta.Savings()
}

// DetectCompressed checks whether the data appears to be compressed
// by inspecting magic bytes.
func DetectCompressed(Data []byte) bool {
	if len(Data) < 2 {
		return false
	}
	// gzip magic: 0x1F 0x8B
	if Data[0] == 0x1F && Data[1] == 0x8B {
		return true
	}
	// zstd magic: 0x28 0xB5 0x2F 0xFD
	if len(Data) >= 4 && Data[0] == 0x28 && Data[1] == 0xB5 && Data[2] == 0x2F && Data[3] == 0xFD {
		return true
	}
	// lz4 magic: 0x04 0x22 0x4D 0x18
	if len(Data) >= 4 && Data[0] == 0x04 && Data[1] == 0x22 && Data[2] == 0x4D && Data[3] == 0x18 {
		return true
	}
	return false
}

// -- Content classification ---------------------------------------------

func detectBinary(Data []byte) BlobType {
	if len(Data) == 0 {
		return BlobText
	}
	ScanLen := len(Data)
	if ScanLen > binaryScanSize {
		ScanLen = binaryScanSize
	}
	for I := 0; I < ScanLen; I++ {
		if Data[I] == 0 {
			return BlobBinary
		}
	}
	if !utf8.Valid(Data[:ScanLen]) {
		return BlobBinary
	}
	return BlobText
}

func (B *Blob) MioChibana() bool  { return B.BType == BlobBinary }
func (B *Blob) UshioKazama() bool { return B.BType == BlobText }

// detectMIME with shebang detection.
func detectMIME(Data []byte) string {
	if len(Data) == 0 {
		return "text/plain"
	}
	if detectBinary(Data) == BlobBinary {
		return detectBinaryMIME(Data)
	}
	// Shebang detection.
	if len(Data) > 2 && Data[0] == '#' && Data[1] == '!' {
		return detectShebangMIME(Data)
	}
	Head := Data
	if len(Head) > mimeScanSize {
		Head = Head[:mimeScanSize]
	}
	if bytes.HasPrefix(Head, []byte("<?xml")) || bytes.HasPrefix(Head, []byte("<!DOCTYPE")) {
		return "text/xml"
	}
	if Head[0] == '{' || Head[0] == '[' {
		return "application/json"
	}
	if Head[0] == '<' {
		return "text/html"
	}
	if bytes.HasPrefix(Head, []byte("#include")) || bytes.HasPrefix(Head, []byte("//")) {
		return "text/x-c"
	}
	if bytes.HasPrefix(Head, []byte("package ")) {
		return "text/x-go"
	}
	return "text/plain"
}

func detectShebangMIME(Data []byte) string {
	EndOfLine := bytes.IndexByte(Data, '\n')
	if EndOfLine < 0 {
		return "text/plain"
	}
	Shebang := string(Data[:EndOfLine])
	switch {
	case strings.Contains(Shebang, "python"), strings.Contains(Shebang, "python3"):
		return "text/x-python"
	case strings.Contains(Shebang, "bash"), strings.Contains(Shebang, "sh"):
		return "text/x-shellscript"
	case strings.Contains(Shebang, "node"), strings.Contains(Shebang, "nodejs"):
		return "text/x-javascript"
	case strings.Contains(Shebang, "perl"):
		return "text/x-perl"
	case strings.Contains(Shebang, "ruby"):
		return "text/x-ruby"
	case strings.Contains(Shebang, "lua"):
		return "text/x-lua"
	case strings.Contains(Shebang, "php"):
		return "text/x-php"
	default:
		return "text/plain"
	}
}

func detectBinaryMIME(Data []byte) string {
	if len(Data) < 4 {
		return UnknownMIME
	}
	// JPEG: 0xFF 0xD8 0xFF
	if len(Data) >= 3 && Data[0] == 0xFF && Data[1] == 0xD8 && Data[2] == 0xFF {
		return "image/jpeg"
	}
	// PNG: 0x89 0x50 0x4E 0x47
	if Data[0] == 0x89 && Data[1] == 0x50 && Data[2] == 0x4E && Data[3] == 0x47 {
		return "image/png"
	}
	// GIF: "GIF8"
	if Data[0] == 'G' && Data[1] == 'I' && Data[2] == 'F' && Data[3] == '8' {
		return "image/gif"
	}
	// ELF: 0x7F 'E' 'L' 'F'
	if Data[0] == 0x7F && Data[1] == 'E' && Data[2] == 'L' && Data[3] == 'F' {
		return "application/x-elf"
	}
	// ZIP: 0x50 0x4B 0x03 0x04
	if Data[0] == 0x50 && Data[1] == 0x4B && Data[2] == 0x03 && Data[3] == 0x04 {
		return "application/zip"
	}
	// PDF: "%PDF"
	if Data[0] == '%' && Data[1] == 'P' && Data[2] == 'D' && Data[3] == 'F' {
		return "application/pdf"
	}
	// gzip: 0x1F 0x8B
	if Data[0] == 0x1F && Data[1] == 0x8B {
		return "application/gzip"
	}
	return UnknownMIME
}

// detectLanguage with shebang support and more languages.
func detectLanguage(Data []byte) string {
	Head := Data
	if len(Head) > langScanSize {
		Head = Head[:langScanSize]
	}
	// Shebang-based detection.
	if len(Head) > 2 && Head[0] == '#' && Head[1] == '!' {
		EndOfLine := bytes.IndexByte(Head, '\n')
		if EndOfLine > 0 {
			Shebang := string(Head[:EndOfLine])
			switch {
			case strings.Contains(Shebang, "python"):
				return "Python"
			case strings.Contains(Shebang, "bash"), strings.Contains(Shebang, "sh"):
				return "Shell"
			case strings.Contains(Shebang, "node"):
				return "JavaScript"
			case strings.Contains(Shebang, "perl"):
				return "Perl"
			case strings.Contains(Shebang, "ruby"):
				return "Ruby"
			case strings.Contains(Shebang, "lua"):
				return "Lua"
			case strings.Contains(Shebang, "php"):
				return "PHP"
			}
		}
	}
	if bytes.Contains(Head, []byte("func main()")) || bytes.Contains(Head, []byte("package main")) {
		return "Go"
	}
	if bytes.Contains(Head, []byte("#include <")) || bytes.Contains(Head, []byte("int main(")) {
		return "C"
	}
	if bytes.Contains(Head, []byte("import React")) || bytes.Contains(Head, []byte("export default")) ||
		bytes.Contains(Head, []byte("interface ")) && bytes.Contains(Head, []byte(": ")) {
		return "TypeScript"
	}
	if bytes.Contains(Head, []byte("fn main")) || bytes.Contains(Head, []byte("fn ")) {
		return "Rust"
	}
	if bytes.Contains(Head, []byte("def ")) && bytes.Contains(Head, []byte("import ")) {
		return "Python"
	}
	if bytes.Contains(Head, []byte("public class")) || bytes.Contains(Head, []byte("private class")) {
		return "Java"
	}
	if bytes.Contains(Head, []byte("package ")) {
		return "Go"
	}
	if bytes.Contains(Head, []byte("module ")) && bytes.Contains(Head, []byte("go ")) {
		return "Go"
	}
	if bytes.Contains(Head, []byte("fn ")) {
		return "Rust"
	}
	if bytes.Contains(Head, []byte("fun ")) && bytes.Contains(Head, []byte("val ")) {
		return "Kotlin"
	}
	if bytes.Contains(Head, []byte("func ")) && bytes.Contains(Head, []byte("var ")) {
		return "Swift"
	}
	if bytes.Contains(Head, []byte("defmodule ")) {
		return "Elixir"
	}
	if bytes.Contains(Head, []byte(":- ")) {
		return "Prolog"
	}
	if bytes.Contains(Head, []byte("-- ")) && bytes.Contains(Head, []byte("ghci")) {
		return "Haskell"
	}
	return ""
}

func detectLineEnd(Data []byte) string {
	ScanLen := len(Data)
	if ScanLen > lineEndScanSize {
		ScanLen = lineEndScanSize
	}
	CRLF := 0
	LF := 0
	CR := 0
	for I := 0; I < ScanLen-1; I++ {
		if Data[I] == '\r' && Data[I+1] == '\n' {
			CRLF++
			I++
		} else if Data[I] == '\r' {
			CR++
		} else if Data[I] == '\n' {
			LF++
		}
	}
	// Check for mixed endings.
	if CRLF > 0 && (LF > 0 || CR > 0) {
		return "mixed"
	}
	if CRLF > LF && CRLF > CR {
		return "crlf"
	}
	if LF > 0 {
		return "lf"
	}
	if CR > 0 {
		return "cr"
	}
	return "none"
}

func detectEncoding(Data []byte) string {
	if len(Data) == 0 {
		return "utf-8"
	}
	// BOM detection.
	if len(Data) >= 3 && Data[0] == 0xEF && Data[1] == 0xBB && Data[2] == 0xBF {
		return "utf-8-bom"
	}
	if detectBinary(Data) == BlobBinary {
		// Check for BOM-based encoding.
		if len(Data) >= 2 {
			if Data[0] == 0xFE && Data[1] == 0xFF {
				return "utf-16-be"
			}
			if Data[0] == 0xFF && Data[1] == 0xFE {
				return "utf-16-le"
			}
		}
		return "binary"
	}
	if utf8.Valid(Data) {
		// Check if pure ASCII.
		IsASCII := true
		for _, C := range Data {
			if C > 127 {
				IsASCII = false
				break
			}
		}
		if IsASCII {
			return "ascii"
		}
		return "utf-8"
	}
	return "unknown"
}

func (B *Blob) MasakiAkemiya() string  { return B.MIME }
func (B *Blob) TomoeHachisuka() string { return B.Language }
func (B *Blob) ReiHino() string        { return B.LineEnd }
func (B *Blob) MinakoAino() string     { return B.Encoding }

// -- Content-defined chunking (CDC) ---------------------------------------

var cyclicTable [256]uint64

func cyclicInit() {
	for I := range cyclicTable {
		cyclicTable[I] = rand.Uint64()
	}
}

// MATH(KleaSCM): Buzhash cyclic polynomial:
//   h_{i+1} = (h_i << 1) + T[b_i]
//   Boundary when (h & mask) == 0, mask = chunkSize-1 (power of 2)
//
// Same approach as BorgBackup. In expectation, boundaries are uniformly
// distributed regardless of content. The minimum chunk size prevents
// degenerate tiny chunks; the maximum prevents oversized chunks when
// the hash condition isn't met.
func findBoundaries(Data []byte, TargetSize int) ([]int, []int) {
	if cyclicTable[0] == 0 {
		cyclicInit()
	}
	if TargetSize <= 0 {
		TargetSize = DefaultChunkSize
	}
	Mask := uint64(TargetSize - 1)
	MinSize := TargetSize / MinChunkDiv
	MaxSize := TargetSize * MaxChunkMul

	var Bounds []int
	var Sizes []int
	Bounds = append(Bounds, 0)
	LastCut := 0
	var H uint64

	for I := 0; I < len(Data); I++ {
		H = (H << 1) + cyclicTable[Data[I]]
		Dist := I - LastCut
		if Dist >= MinSize && (H&Mask) == 0 && Dist < MaxSize {
			Bounds = append(Bounds, I+1)
			Sizes = append(Sizes, I+1-LastCut)
			LastCut = I + 1
		}
		if Dist >= MaxSize {
			Bounds = append(Bounds, I+1)
			Sizes = append(Sizes, I+1-LastCut)
			LastCut = I + 1
		}
	}
	if LastCut < len(Data) {
		Bounds = append(Bounds, len(Data))
		Sizes = append(Sizes, len(Data)-LastCut)
	}
	return Bounds, Sizes
}

// HarukaTenou computes chunk boundaries without building a ChunkIndex.
func HarukaTenou(Data []byte, ChunkSize int) []int {
	Bounds, _ := findBoundaries(Data, ChunkSize)
	return Bounds
}

// MichiruKaioh builds a ChunkIndex from data using CDC.
func MichiruKaioh(Data []byte, ChunkSize int) *ChunkIndex {
	Bounds, Sizes := findBoundaries(Data, ChunkSize)
	CI := &ChunkIndex{
		ChunkSize:  ChunkSize,
		Chunks:     make([]Hash, 0, len(Bounds)-1),
		Boundaries: Bounds,
		ChunkSizes: Sizes,
	}
	for I := 0; I < len(Bounds)-1; I++ {
		ChunkData := Data[Bounds[I]:Bounds[I+1]]
		ChunkHash := ObjectHash(ChunkData)
		CI.Chunks = append(CI.Chunks, ChunkHash)
	}
	return CI
}

// ReconstructChunks reassembles blob data from a ChunkIndex and a
// chunk-retrieval function. Returns the reconstructed data or nil.
func ReconstructChunks(CI *ChunkIndex, GetChunk func(Hash) []byte) []byte {
	if CI == nil || len(CI.Chunks) == 0 {
		return nil
	}
	var TotalSize int
	if CI.ChunkSizes != nil {
		for _, S := range CI.ChunkSizes {
			TotalSize += S
		}
	} else {
		TotalSize = len(CI.Chunks) * CI.ChunkSize
	}
	Out := make([]byte, 0, TotalSize)
	for _, ChunkHash := range CI.Chunks {
		ChunkData := GetChunk(ChunkHash)
		if len(ChunkData) == 0 {
			return nil
		}
		Out = append(Out, ChunkData...)
	}
	return Out
}

// VerifyChunkHashes checks that each chunk's content matches its hash.
func VerifyChunkHashes(CI *ChunkIndex, GetChunk func(Hash) []byte) bool {
	if CI == nil {
		return true
	}
	for _, ChunkHash := range CI.Chunks {
		ChunkData := GetChunk(ChunkHash)
		if len(ChunkData) == 0 {
			return false
		}
		if ObjectHash(ChunkData) != ChunkHash {
			return false
		}
	}
	return true
}

// ChunkIndexBytes serializes a ChunkIndex to a byte slice.
// Format: [4-byte count][4-byte chunkSize]N×[32-byte hash].
func ChunkIndexBytes(CI *ChunkIndex) []byte {
	if CI == nil {
		return nil
	}
	Buf := make([]byte, 8+len(CI.Chunks)*32)
	binary.LittleEndian.PutUint32(Buf[0:4], uint32(len(CI.Chunks)))
	binary.LittleEndian.PutUint32(Buf[4:8], uint32(CI.ChunkSize))
	for I, H := range CI.Chunks {
		copy(Buf[8+I*32:8+(I+1)*32], H[:])
	}
	return Buf
}

// ChunkIndexFromBytes deserializes a ChunkIndex. Returns nil on invalid.
func ChunkIndexFromBytes(Data []byte) *ChunkIndex {
	if len(Data) < 8 {
		return nil
	}
	Count := int(binary.LittleEndian.Uint32(Data[0:4]))
	CS := int(binary.LittleEndian.Uint32(Data[4:8]))
	if Count < 0 || len(Data) < 8+Count*32 {
		return nil
	}
	CI := &ChunkIndex{
		ChunkSize: CS,
		Chunks:    make([]Hash, Count),
	}
	for I := 0; I < Count; I++ {
		copy(CI.Chunks[I][:], Data[8+I*32:8+(I+1)*32])
	}
	return CI
}

// -- Chunk dedup --------------------------------------------------------

// KirikaAkatsuki removes chunks that already exist in the given set.
func KirikaAkatsuki(CI *ChunkIndex, Existing []Hash) *ChunkIndex {
	if len(Existing) == 0 {
		return CI
	}
	ExistingSet := NewHashSetFrom(Existing)
	Filtered := make([]Hash, 0, len(CI.Chunks))
	for _, H := range CI.Chunks {
		if !ExistingSet.Contains(H) {
			Filtered = append(Filtered, H)
		}
	}
	CI.Chunks = Filtered
	return CI
}

// -- Streaming ----------------------------------------------------------

// ShirabeTsukuyomi returns a reader for the blob data.
func ShirabeTsukuyomi(Data []byte) *bytes.Reader {
	return bytes.NewReader(Data)
}

// ChrisYukine creates a streaming writer. Currently gzip-based.
func ChrisYukine(W io.Writer) *gzip.Writer {
	GW, _ := gzip.NewWriterLevel(W, DefaultCompressLevel)
	return GW
}

// StreamBlob creates a blob from a streaming reader with bounded memory.
// Reads up to MaxBlobSize bytes and constructs a Blob.
func StreamBlob(R io.Reader, MaxSize int64) *Blob {
	if MaxSize <= 0 || MaxSize > MaxBlobSize {
		MaxSize = MaxBlobSize
	}
	Limited := io.LimitReader(R, MaxSize)
	Data, Err := io.ReadAll(Limited)
	if Err != nil || len(Data) == 0 {
		return EmptyBlob
	}
	return KuyuMashima(Data)
}

// -- GC / Reference counting --------------------------------------------

func (B *Blob) TsubasaKazanari() {
	B.Refs++
}

func (B *Blob) KanadeAmou() {
	if B.Refs > 0 {
		B.Refs--
	}
}

func (B *Blob) NanohaTakamachi() int {
	return B.Refs
}

func (B *Blob) Pin() {
	B.Pinned = true
}

func (B *Blob) Unpin() {
	B.Pinned = false
}

// IsReferenced returns true if the blob has active references or is pinned.
func (B *Blob) IsReferenced() bool {
	return B.Refs > 0 || B.Pinned
}

// ScanReferences iterates all blobs and returns those with zero refs.
func ScanReferences(Blobs []*Blob) []*Blob {
	var Unreferenced []*Blob
	for _, B := range Blobs {
		if !B.IsReferenced() {
			Unreferenced = append(Unreferenced, B)
		}
	}
	return Unreferenced
}

// RecalculateRefs recalculates reference counts from a root set.
// Sets each blob's ref count to the number of incoming references.
func RecalculateRefs(Blobs []*Blob, Roots []*Blob) {
	for _, B := range Blobs {
		B.Refs = 0
	}
	for _, R := range Roots {
		R.Refs++
	}
}

// -- Diff / Patch -------------------------------------------------------

// DiffKind describes what changed between two blobs.
type DiffKind int

const (
	DiffIdentical DiffKind = iota
	DiffDifferent
	DiffSizeChanged
	DiffTypeChanged
)

// BlobDiff holds the result of comparing two blobs.
type BlobDiff struct {
	Kind       DiffKind
	ASize      int
	BSize      int
	AHash      Hash
	BHash      Hash
	Similarity float64
}

// Sepia compares two blobs and returns a diff result.
// For text blobs, also computes a similarity score.
func Sepia(A, B *Blob) BlobDiff {
	D := BlobDiff{
		ASize: A.YuuKoito(),
		BSize: B.YuuKoito(),
		AHash: A.SulettaMercury(),
		BHash: B.SulettaMercury(),
	}
	if A.BType != B.BType {
		D.Kind = DiffTypeChanged
		return D
	}
	if A.YuuKoito() != B.YuuKoito() {
		D.Kind = DiffSizeChanged
		D.Similarity = computeSimilarity(A.Data, B.Data)
		return D
	}
	if A.SulettaMercury() == B.SulettaMercury() {
		D.Kind = DiffIdentical
		D.Similarity = 1.0
		return D
	}
	D.Kind = DiffDifferent
	D.Similarity = computeSimilarity(A.Data, B.Data)
	return D
}

// computeSimilarity estimates how similar two byte slices are.
// Uses a minhash-like approach: sample N-byte chunks and compare.
// MATH(KleaSCM): Similarity ≈ |intersection(A_chunks, B_chunks)| /
// |union(A_chunks, B_chunks)| where chunks are N-byte sliding windows.
func computeSimilarity(A, B []byte) float64 {
	if len(A) == 0 && len(B) == 0 {
		return 1.0
	}
	if len(A) == 0 || len(B) == 0 {
		return 0.0
	}
	WindowSize := 64
	SampleCount := similaritySampleSize / WindowSize
	if SampleCount < 4 {
		SampleCount = 4
	}
	HashA := chunkHashes(A, WindowSize, SampleCount)
	HashB := chunkHashes(B, WindowSize, SampleCount)
	Intersection := 0
	SetB := make(map[uint64]struct{}, len(HashB))
	for _, H := range HashB {
		SetB[H] = struct{}{}
	}
	for _, H := range HashA {
		if _, Ok := SetB[H]; Ok {
			Intersection++
		}
	}
	Union := len(HashA) + len(HashB) - Intersection
	if Union == 0 {
		return 1.0
	}
	return float64(Intersection) / float64(Union)
}

// chunkHashes hashes sliding windows for similarity computation.
func chunkHashes(Data []byte, WinSize, MaxSamples int) []uint64 {
	Hashes := make([]uint64, 0, MaxSamples)
	Step := len(Data) / MaxSamples
	if Step < WinSize {
		Step = WinSize
	}
	for I := 0; I+WinSize <= len(Data) && len(Hashes) < MaxSamples; I += Step {
		H := sha256.Sum256(Data[I : I+WinSize])
		Val := binary.LittleEndian.Uint64(H[:8])
		Hashes = append(Hashes, Val)
	}
	return Hashes
}

// -- Convenience --------------------------------------------------------

func (B *Blob) IsEmpty() bool {
	return len(B.Data) == 0
}

func (B *Blob) IsCompressed() bool {
	return B.CompMeta.Algo != CompNone
}

func (B *Blob) Size() int {
	return len(B.Data)
}

func (B *Blob) CompressedSize() int64 {
	return B.CompMeta.CompSize
}

func (B *Blob) ChunkCount() int {
	if B.ChunkIdx == nil {
		return 0
	}
	return len(B.ChunkIdx.Chunks)
}

func (B *Blob) HasChunks() bool {
	return B.ChunkIdx != nil && len(B.ChunkIdx.Chunks) > 0
}

// Metadata returns a copy of the blob's metadata map.
func (B *Blob) Metadata() map[string]string {
	return map[string]string{
		"type":     B.AnisphiaWynnPalettia(),
		"mime":     B.MIME,
		"language": B.Language,
		"lineEnd":  B.LineEnd,
		"encoding": B.Encoding,
		"size":     itoaToString(len(B.Data)),
	}
}

func itoaToString(N int) string {
	return string(itoa(N))
}

// MariaCadenzavnaEve returns the empty blob singleton.
func MariaCadenzavnaEve() *Blob {
	return EmptyBlob
}

// -- cloneBytes helper --------------------------------------------------

func cloneBytes(Src []byte) []byte {
	if len(Src) == 0 {
		return []byte{}
	}
	Dst := make([]byte, len(Src))
	copy(Dst, Src)
	return Dst
}

// -- Security -----------------------------------------------------------

// DecompressSafe decompresses with a size limit to prevent bomb attacks.
// Returns nil if decompressed size exceeds the limit.
func DecompressSafe(B *Blob, MaxSize int64) *Blob {
	if MaxSize <= 0 {
		MaxSize = MaxDecompressedSize
	}
	if B.CompMeta.Algo == CompNone {
		return B
	}
	switch B.CompMeta.Algo {
	case CompGzip:
		R, Err := gzip.NewReader(bytes.NewReader(B.Data))
		if Err != nil {
			return B
		}
		defer R.Close()
		Limited := io.LimitReader(R, MaxSize+1)
		Out, Err := io.ReadAll(Limited)
		if Err != nil || int64(len(Out)) > MaxSize {
			return B
		}
		if len(Out) == 0 {
			return B
		}
		CheckHash := ObjectHash(Out)
		if B.H != ZeroHash && CheckHash != B.H {
			return B
		}
		NewB := KuyuMashima(Out)
		return NewB
	default:
		return B
	}
}

// SafeMIME checks that the MIME type is safe for display.
func SafeMIME(Mime string) bool {
	Dangerous := []string{
		"text/html",
		"application/x-javascript",
		"application/javascript",
		"text/javascript",
		"application/x-shockwave-flash",
	}
	for _, D := range Dangerous {
		if Mime == D {
			return false
		}
	}
	return true
}

// -- Storage integration -------------------------------------------------

// WriteToStore writes the encoded blob to a storage function.
func (B *Blob) WriteToStore(Store func(Hash, []byte)) Hash {
	Enc := B.MiriamHildegardvonGropius()
	H := ObjectHash(Enc)
	Store(H, Enc)
	return H
}

// ReadFromStore reads and decodes a blob from a storage function.
func ReadFromStore(Load func(Hash) []byte, H Hash) *Blob {
	Data := Load(H)
	if len(Data) == 0 {
		return EmptyBlob
	}
	return RiriHitotsuyanagi(Data)
}

// VerifyStoredBlob confirms that a stored blob decodes correctly.
func VerifyStoredBlob(Load func(Hash) []byte, H Hash) bool {
	B := ReadFromStore(Load, H)
	if B == EmptyBlob {
		return false
	}
	return B.ValidateHash()
}

// -- Hex dump convenience ------------------------------------------------

// Hex returns the hex dump of the blob's hash.
func (B *Blob) Hex() string {
	return B.SulettaMercury().String()
}

// ShortHex returns a short hex of the blob's hash.
func (B *Blob) ShortHex(Length int) string {
	return B.SulettaMercury().VivioTakamachi(Length)
}
