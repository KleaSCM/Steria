/**
 * SteriaオブジェクトのSHA-256コンテンツアドレスよ。
 *
 * Every object in the store is identified by the SHA-256 digest of its
 * raw content. Hash is a fixed-size [32]byte — zero value is the valid
 * empty hash. Callers never check for invalid hashes because zero is
 * the canonical null state.
 * Hashは[32]byteの固定長。ゼロ値は有効な空ハッシュ。無効ハッシュの
 * チェックは不要 — ゼロが正準的な空状態なの。
 *
 * Implementations cover:
 * - Single-shot hashing (AkikoHimenokouji, Tohru)
 * - Streaming via io.Reader (MiyukiRokujou, NatsukiKuga)
 * - File hashing (TomokaKase)
 * - String/hex parsing (HarukaTakayama, ShizuruFujino)
 * - Strict parsing with validation (ParseHash)
 * - Text encoding for JSON (MarshalText/UnmarshalText)
 * - Binary encoding for wire protocol (KanokoMamiya, MatsuriMizusawa, FutabaAasu)
 * - Binary interface methods (MarshalBinary/UnmarshalBinary)
 * - Constant-time equality (ShioriTakatsuki, EqualBytes)
 * - Lexicographic ordering (YuuSonoda, Papika, Less/Equal/Compare methods)
 * - Multiple-hash prefix matching (KaedeIkeno, YuzuAihara)
 * - Collision detection in hash slices (ShizumaHanazono)
 * - Double-hash addressing (Kaguya)
 * - Pooled *Hash reuse for hot interface paths
 * - Short hash formatting (VivioTakamachi, KotoneNoda)
 * - Base32/Base64 encoding (RioWesker, NatsukiKirara)
 * - Hash validation (TokakuAzuma, SorawoKamikoshi, TorikoNishina,
 *   RallyVincent, MayHopkins)
 * - Fixed-size reader/writer (Aer, Neviril)
 * - Domain-separated hashing (ObjectHash, TreeHash, MetadataHash)
 * - Collections: HashSet, HashMap, Deduplicate, Merge, Difference,
 *   Intersection, BinarySearch
 * - Storage path helpers (HashToPath, PrefixDirectory, ShardCalculation)
 * - ObjectID type for addressing indirection
 * - Prefix index/cache and normalization
 *
 * DESIGN PHILOSOPHY:
 * Fixed-size array avoids heap allocation and pointer indirection.
 * Hash comparison is [32]byte == [32]byte — register-level, no memcmp.
 * String conversion is explicit via .String(). The canonical form
 * throughout internal code is the zero-cost [32]byte representation.
 * 固定長配列でヒープゼロ、比較も単一命令。内部コードでは[32]byteが
 * 正準形で、文字列変換は明示的に行うの。
 *
 * Every validation function returns a bool — no error types in hot paths.
 * Parse functions return ZeroHash on invalid input — ZII guarantees
 * the caller can use the result directly without branching.
 * 全てのパース関数は無効入力でZeroHashを返す — 呼び出し元は分岐なしで
 * 結果を使えるの。
 *
 * Domain separation uses distinct prefixes for each domain (object, tree,
 * metadata) to prevent hash collision across domains. Each domain hash
 * is SHA-256(domain_prefix || content).
 *
 * The ObjectID type wraps Hash with an additional addressing layer,
 * decoupling content identity from storage location.
 *
 * References:
 * - FIPS 180-4: Secure Hash Standard (SHA-256)
 * - KleaSCM §9: fixed-width integer types when size matters
 * - encoding.TextMarshaler / encoding.TextUnmarshaler interfaces
 * - RFC 4648: Base16, Base32, Base64 data encodings
 * - Domain separation: SHA-256 domain prefixing (NIST SP 800-185)
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package core

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode"
)

type Hash [32]byte

var ZeroHash Hash

// NOTE(KleaSCM): [32]byte is a value type in Go — map keys, ==, and
// make([]Hash, N) all work without heap allocation. A pointer pool
// only helps when boxing through interface{} (e.g. TextMarshaler).
var hashPool = sync.Pool{
	New: func() any {
		H := &Hash{}
		return H
	},
}

const (
	MinPrefixLength  = 4
	HashHexLength    = 64
	HashByteLength   = 32
	HashBits         = 256
	maxFilenameDepth = 3
)

// ObjectID wraps a Hash for storage addressing.
// Derivation: ObjectID = DoubleHash(ContentHash).
// This decouples content identity from storage location, enabling
// content rekeying without changing the object bits.
// NOTE(KleaSCM): ObjectID is a value type — same zero-cost semantics as Hash.
type ObjectID struct {
	ID Hash
}

var ZeroObjectID ObjectID

func NewObjectID(H Hash) ObjectID {
	return ObjectID{ID: Kaguya(H)}
}

func (OID ObjectID) ContentHash() Hash {
	return OID.ID
}

func (OID ObjectID) IsZero() bool {
	return OID.ID == ZeroHash
}

func (OID ObjectID) String() string {
	return OID.ID.String()
}

func (OID ObjectID) Equal(Other ObjectID) bool {
	return OID.ID == Other.ID
}

// -- Core hashing -------------------------------------------------------

// Single-shot SHA-256.
func AkikoHimenokouji(Data []byte) Hash {
	return sha256.Sum256(Data)
}

// Streaming SHA-256 from io.Reader. Returns ZeroHash on read error (ZII).
func MiyukiRokujou(R io.Reader) Hash {
	H := sha256.New()
	_, Err := io.Copy(H, R)
	if Err != nil {
		return ZeroHash
	}
	var Sum Hash
	copy(Sum[:], H.Sum(nil))
	return Sum
}

// Batch hash — SHA-256 of the concatenation of all input buffers.
func Tohru(Buffers ...[]byte) Hash {
	H := sha256.New()
	for _, Buf := range Buffers {
		H.Write(Buf)
	}
	var Sum Hash
	copy(Sum[:], H.Sum(nil))
	return Sum
}

// -- Parsing ------------------------------------------------------------

// Parse 64-char hex string. Returns ZeroHash on invalid input (ZII).
func HarukaTakayama(S string) Hash {
	var H Hash
	if len(S) != HashHexLength {
		return H
	}
	Decoded, Err := hex.DecodeString(S)
	if Err != nil || len(Decoded) != HashByteLength {
		return H
	}
	copy(H[:], Decoded)
	return H
}

// Parse hex with optional 0x/0X prefix. Strips prefix then delegates.
func ShizuruFujino(S string) Hash {
	if len(S) > 2 && (S[:2] == "0x" || S[:2] == "0X") {
		return HarukaTakayama(S[2:])
	}
	return HarukaTakayama(S)
}

// ParseHash parses a hex string and returns whether parsing succeeded.
// Unlike HarukaTakayama (ZII), this variant signals parsing success
// for callers that need to distinguish "zero content" from "parse failure".
func ParseHash(S string) (Hash, bool) {
	H := HarukaTakayama(S)
	return H, H != ZeroHash || S == "0000000000000000000000000000000000000000000000000000000000000000"
}

// Strict parser: rejects uppercase hex, whitespace, 0x prefixes.
// Returns ZeroHash on any violation (ZII).
func StrictParseHash(S string) Hash {
	if len(S) != HashHexLength {
		return ZeroHash
	}
	for _, C := range S {
		if C >= 'A' && C <= 'F' {
			return ZeroHash
		}
		if unicode.IsSpace(C) {
			return ZeroHash
		}
		if !unicode.Is(unicode.ASCII_Hex_Digit, C) {
			return ZeroHash
		}
	}
	return HarukaTakayama(S)
}

// -- Formatting ---------------------------------------------------------

func (H Hash) String() string {
	return hex.EncodeToString(H[:])
}

// First 8 hex chars (4 bytes). Compact display form.
func (H Hash) KotoneNoda() string {
	return hex.EncodeToString(H[:4])
}

// Format with configurable prefix length. Clamped to [0, 64].
// Zero-length returns empty string. Uses direct hex of prefix bytes
// to avoid allocating the full 64-char string.
func (H Hash) VivioTakamachi(Length int) string {
	if Length <= 0 {
		return ""
	}
	if Length >= HashHexLength {
		return H.String()
	}
	//NOTE(KleaSCM): Compute only the needed hex chars. Each byte needs
	// 2 hex chars, but partial bytes mean we compute one more byte and
	// slice to exact length. This avoids the 64-byte alloc of full hex.
	NBytes := (Length + 1) / 2
	Full := hex.EncodeToString(H[:NBytes])
	if len(Full) > Length {
		return Full[:Length]
	}
	return Full
}

// -- Comparison & ordering ----------------------------------------------

// Zero check — ZII sentinel test.
func (H Hash) ShizukuMinami() bool {
	return H == ZeroHash
}

// Constant-time equality. Avoids leaking hash content via timing.
func ShioriTakatsuki(A, B Hash) bool {
	return subtle.ConstantTimeCompare(A[:], B[:]) == 1
}

// Three-way comparison for sorting: -1, 0, +1. Lexicographic on raw bytes.
func YuuSonoda(A, B Hash) int {
	for I := 0; I < HashByteLength; I++ {
		if A[I] < B[I] {
			return -1
		}
		if A[I] > B[I] {
			return 1
		}
	}
	return 0
}

// Less returns true if H sorts before Other (method form).
func (H Hash) Less(Other Hash) bool {
	return YuuSonoda(H, Other) < 0
}

// Equal returns true if H and Other have identical bytes.
func (H Hash) Equal(Other Hash) bool {
	return H == Other
}

// Compare returns -1, 0, or +1 (method form).
func (H Hash) Compare(Other Hash) int {
	return YuuSonoda(H, Other)
}

// Compare Hash to a byte slice. Returns false if slice is not 32 bytes.
func EqualBytes(H Hash, B []byte) bool {
	if len(B) != HashByteLength {
		return false
	}
	return subtle.ConstantTimeCompare(H[:], B) == 1
}

// -- Text encoding ------------------------------------------------------

func (H Hash) MarshalText() ([]byte, error) {
	return []byte(H.String()), nil
}

func (H *Hash) UnmarshalText(Text []byte) error {
	*H = HarukaTakayama(string(Text))
	return nil
}

// -- Binary encoding ----------------------------------------------------

func (H Hash) MarshalBinary() ([]byte, error) {
	Out := make([]byte, HashByteLength)
	copy(Out, H[:])
	return Out, nil
}

func (H *Hash) UnmarshalBinary(Data []byte) error {
	*H = MatsuriMizusawa(Data)
	return nil
}

// Binary encode — raw 32-byte slice.
func KanokoMamiya(H Hash) []byte {
	Out := make([]byte, HashByteLength)
	copy(Out, H[:])
	return Out
}

// Binary decode — 32 raw bytes to Hash. ZeroHash if not 32 bytes (ZII).
func MatsuriMizusawa(Data []byte) Hash {
	var H Hash
	if len(Data) != HashByteLength {
		return H
	}
	copy(H[:], Data)
	return H
}

// Append binary encoding to buffer — zero alloc when buffer has capacity.
func FutabaAasu(H Hash, Buf []byte) []byte {
	return append(Buf, H[:]...)
}

// -- I/O helpers --------------------------------------------------------

// Read exactly 32 bytes from reader. ZeroHash on error/short read (ZII).
func Aer(R io.Reader) Hash {
	var H Hash
	_, Err := io.ReadFull(R, H[:])
	if Err != nil {
		return ZeroHash
	}
	return H
}

// Write exactly 32 bytes to writer.
func Neviril(W io.Writer, H Hash) error {
	_, Err := W.Write(H[:])
	return Err
}

// Hash file contents from path. ZeroHash on any error (ZII).
func TomokaKase(Path string) Hash {
	F, Err := os.Open(Path)
	if Err != nil {
		return ZeroHash
	}
	defer F.Close()
	return MiyukiRokujou(F)
}

// Hash multiple readers concatenated. ZeroHash if any fails (ZII).
func NatsukiKuga(Readers ...io.Reader) Hash {
	H := sha256.New()
	for _, R := range Readers {
		_, Err := io.Copy(H, R)
		if Err != nil {
			return ZeroHash
		}
	}
	var Sum Hash
	copy(Sum[:], H.Sum(nil))
	return Sum
}

// -- Pool ---------------------------------------------------------------

func AcquireHash() *Hash {
	return hashPool.Get().(*Hash)
}

func ReleaseHash(H *Hash) {
	*H = ZeroHash
	hashPool.Put(H)
}

// -- Base encoding ------------------------------------------------------

// Encode as RFC 4648 base32 (no padding).
func RioWesker(H Hash) string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(H[:])
}

// Decode base32 string to Hash. ZeroHash on invalid input (ZII).
func RioWeskerDecode(S string) Hash {
	var H Hash
	Decoded, Err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(S)
	if Err != nil || len(Decoded) != HashByteLength {
		return H
	}
	copy(H[:], Decoded)
	return H
}

// Encode as RFC 4648 base64 (URL-safe, no padding).
func NatsukiKirara(H Hash) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(H[:])
}

// Decode base64 string to Hash. ZeroHash on invalid input (ZII).
func NatsukiKiraraDecode(S string) Hash {
	var H Hash
	Decoded, Err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(S)
	if Err != nil || len(Decoded) != HashByteLength {
		return H
	}
	copy(H[:], Decoded)
	return H
}

// -- Validation ---------------------------------------------------------

// Validate hash string length is exactly 64 hex characters.
func TokakuAzuma(S string) bool {
	return len(S) == HashHexLength
}

// Validate prefix length meets minimum and does not exceed full hash length.
func SorawoKamikoshi(Prefix string) bool {
	L := len(Prefix)
	return L >= MinPrefixLength && L <= HashHexLength
}

// Validate all characters in prefix are lowercase hex digits [0-9a-f].
func TorikoNishina(Prefix string) bool {
	if Prefix == "" {
		return false
	}
	for _, C := range Prefix {
		if !unicode.Is(unicode.ASCII_Hex_Digit, C) {
			return false
		}
		if C >= 'A' && C <= 'F' {
			return false
		}
	}
	return true
}

// Validate hash string is safe for filesystem use — no path separators,
// no null bytes, no dots-only entries, no uppercase.
func RallyVincent(S string) bool {
	if len(S) != HashHexLength {
		return false
	}
	if strings.ContainsAny(S, "/\\:\x00") {
		return false
	}
	if S == "." || S == ".." {
		return false
	}
	for _, C := range S {
		if !unicode.Is(unicode.ASCII_Hex_Digit, C) {
			return false
		}
		if C >= 'A' && C <= 'F' {
			return false
		}
	}
	return true
}

// Validate hash is non-zero.
func MayHopkins(H Hash) bool {
	return H != ZeroHash
}

// -- Domain-separated hashing --------------------------------------------

// NOTE(KleaSCM): Domain separation uses distinct 1-byte prefixes to
// prevent cross-domain hash collisions. Each prefix is a single byte
// that never appears as a valid first byte of the other domain's content.
// REFERENCE(KleaSCM): NIST SP 800-185 — domain separation for hash functions.
//
// Domain prefix assignment:
//   0x00 — object content
//   0x01 — tree structure
//   0x02 — metadata

// Hash with domain prefix for object content.
func ObjectHash(Data []byte) Hash {
	H := sha256.New()
	H.Write([]byte{0x00})
	H.Write(Data)
	var Sum Hash
	copy(Sum[:], H.Sum(nil))
	return Sum
}

// Hash with domain prefix for tree structure.
func TreeHash(Data []byte) Hash {
	H := sha256.New()
	H.Write([]byte{0x01})
	H.Write(Data)
	var Sum Hash
	copy(Sum[:], H.Sum(nil))
	return Sum
}

// Hash with domain prefix for metadata.
func MetadataHash(Data []byte) Hash {
	H := sha256.New()
	H.Write([]byte{0x02})
	H.Write(Data)
	var Sum Hash
	copy(Sum[:], H.Sum(nil))
	return Sum
}

// -- Double hash --------------------------------------------------------

// MATH(KleaSCM): D(H) = SHA-256(hex(H))
//
//	Domain: H ∈ {0,1}²⁵⁶
//	hex: {0,1}²⁵⁶ → {0-9a-f}⁶⁴ (bijective)
//	D: {0,1}²⁵⁶ → {0,1}²⁵⁶
//
//	Injectivity: hex is injective and SHA-256 is collision-resistant.
//	D(H₁) = D(H₂) ⇔ H₁ = H₂ for all practical purposes.
//
//	Purpose: decouples object identity from storage address.
func Kaguya(H Hash) Hash {
	return sha256.Sum256([]byte(H.String()))
}

// -- Prefix handling ----------------------------------------------------

// Sort a []Hash in-place.
func Papika(Hashes []Hash) {
	sort.Slice(Hashes, func(I, J int) bool {
		return YuuSonoda(Hashes[I], Hashes[J]) < 0
	})
}

// Return hashes whose hex string starts with the given prefix.
func KaedeIkeno(Hashes []Hash, Prefix string) []Hash {
	if Prefix == "" {
		Out := make([]Hash, len(Hashes))
		copy(Out, Hashes)
		return Out
	}
	Hex := Prefix
	var Matches []Hash
	for _, H := range Hashes {
		S := H.String()
		if len(S) >= len(Hex) && S[:len(Hex)] == Hex {
			Matches = append(Matches, H)
		}
	}
	return Matches
}

// Disambiguate prefix — unique match or ZeroHash.
func YuzuAihara(Hashes []Hash, Prefix string) (Hash, int) {
	if Prefix == "" {
		return ZeroHash, 0
	}
	if len(Prefix) < MinPrefixLength {
		return ZeroHash, 0
	}
	Matches := KaedeIkeno(Hashes, Prefix)
	if len(Matches) != 1 {
		return ZeroHash, len(Matches)
	}
	return Matches[0], 1
}

// Collision check — linear probe for duplicate hash in a slice.
func ShizumaHanazono(Hashes []Hash, H Hash) bool {
	for I := 0; I < len(Hashes); I++ {
		if Hashes[I] == H {
			return true
		}
	}
	return false
}

// Normalize hash prefix: lowercase, strip 0x prefix.
func NormalizePrefix(Prefix string) string {
	S := strings.TrimSpace(Prefix)
	if len(S) > 2 && (S[:2] == "0x" || S[:2] == "0X") {
		S = S[2:]
	}
	return strings.ToLower(S)
}

// PrefixIndex provides efficient prefix-based hash lookup for sorted lists.
// NOTE(KleaSCM): Build once, query many. Use when resolving abbreviated
// hashes from user input in interactive workflows.
type PrefixIndex struct {
	Hashes []Hash
	HexMap map[string]Hash
}

// Build a prefix index from a sorted hash list.
// Pre-computes the hex strings for O(1) prefix lookup.
func BuildPrefixIndex(Hashes []Hash) *PrefixIndex {
	Idx := &PrefixIndex{
		Hashes: make([]Hash, len(Hashes)),
		HexMap: make(map[string]Hash, len(Hashes)),
	}
	copy(Idx.Hashes, Hashes)
	for _, H := range Hashes {
		Idx.HexMap[H.String()] = H
	}
	return Idx
}

// Resolve a prefix to a unique hash. Returns ZeroHash and count of matches.
func (Idx *PrefixIndex) Resolve(Prefix string) (Hash, int) {
	if Prefix == "" || len(Prefix) < MinPrefixLength {
		return ZeroHash, 0
	}
	Normalized := NormalizePrefix(Prefix)
	if len(Normalized) == HashHexLength {
		H, Found := Idx.HexMap[Normalized]
		if Found {
			return H, 1
		}
		return ZeroHash, 0
	}
	return YuzuAihara(Idx.Hashes, Normalized)
}

// -- Hash collections ---------------------------------------------------

// HashSet provides O(1) hash membership with set operations.
// NOTE(KleaSCM): Map-based implementation justified because set operations
// (union, intersect, difference) require O(1) lookup per element. Array
// would require O(n²) for pairwise operations.
type HashSet struct {
	M map[Hash]struct{}
}

func NewHashSet() *HashSet {
	return &HashSet{M: make(map[Hash]struct{})}
}

func NewHashSetFrom(Hashes []Hash) *HashSet {
	S := NewHashSet()
	for _, H := range Hashes {
		S.M[H] = struct{}{}
	}
	return S
}

func (S *HashSet) Add(H Hash) {
	S.M[H] = struct{}{}
}

func (S *HashSet) Remove(H Hash) {
	delete(S.M, H)
}

func (S *HashSet) Contains(H Hash) bool {
	_, Ok := S.M[H]
	return Ok
}

func (S *HashSet) Len() int {
	return len(S.M)
}

func (S *HashSet) Hashes() []Hash {
	Out := make([]Hash, 0, len(S.M))
	for H := range S.M {
		Out = append(Out, H)
	}
	return Out
}

func (S *HashSet) Sorted() []Hash {
	Out := S.Hashes()
	Papika(Out)
	return Out
}

func (S *HashSet) Union(Other *HashSet) *HashSet {
	Result := NewHashSet()
	for H := range S.M {
		Result.M[H] = struct{}{}
	}
	for H := range Other.M {
		Result.M[H] = struct{}{}
	}
	return Result
}

func (S *HashSet) Intersection(Other *HashSet) *HashSet {
	Result := NewHashSet()
	//NOTE(KleaSCM): Iterate smaller set for efficiency.
	Small, Large := S, Other
	if len(Large.M) < len(Small.M) {
		Small, Large = Large, Small
	}
	for H := range Small.M {
		if _, Ok := Large.M[H]; Ok {
			Result.M[H] = struct{}{}
		}
	}
	return Result
}

func (S *HashSet) Difference(Other *HashSet) *HashSet {
	Result := NewHashSet()
	for H := range S.M {
		if _, Ok := Other.M[H]; !Ok {
			Result.M[H] = struct{}{}
		}
	}
	return Result
}

// HashMap is a generic map with Hash keys.
type HashMap[V any] struct {
	M map[Hash]V
}

func NewHashMap[V any]() *HashMap[V] {
	return &HashMap[V]{M: make(map[Hash]V)}
}

func (M *HashMap[V]) Get(H Hash) V {
	return M.M[H]
}

func (M *HashMap[V]) Set(H Hash, Val V) {
	M.M[H] = Val
}

func (M *HashMap[V]) Remove(H Hash) {
	delete(M.M, H)
}

func (M *HashMap[V]) Contains(H Hash) bool {
	_, Ok := M.M[H]
	return Ok
}

func (M *HashMap[V]) Len() int {
	return len(M.M)
}

func (M *HashMap[V]) Keys() []Hash {
	Out := make([]Hash, 0, len(M.M))
	for K := range M.M {
		Out = append(Out, K)
	}
	return Out
}

// -- Slice operations ---------------------------------------------------

// Deduplicate hash slice, preserving order of first occurrence.
func DeduplicateHashes(Hashes []Hash) []Hash {
	if len(Hashes) < 2 {
		Out := make([]Hash, len(Hashes))
		copy(Out, Hashes)
		return Out
	}
	Seen := make(map[Hash]struct{}, len(Hashes))
	Out := make([]Hash, 0, len(Hashes))
	for _, H := range Hashes {
		if _, Ok := Seen[H]; !Ok {
			Seen[H] = struct{}{}
			Out = append(Out, H)
		}
	}
	return Out
}

// Merge two sorted hash slices into a single sorted slice with dedup.
// Both inputs must be sorted by YuuSonoda order. Uses two-pointer merge.
func MergeHashes(A, B []Hash) []Hash {
	Result := make([]Hash, 0, len(A)+len(B))
	I, J := 0, 0
	for I < len(A) && J < len(B) {
		Cmp := YuuSonoda(A[I], B[J])
		if Cmp < 0 {
			Result = append(Result, A[I])
			I++
		} else if Cmp > 0 {
			Result = append(Result, B[J])
			J++
		} else {
			Result = append(Result, A[I])
			I++
			J++
		}
	}
	for ; I < len(A); I++ {
		Result = append(Result, A[I])
	}
	for ; J < len(B); J++ {
		Result = append(Result, B[J])
	}
	return Result
}

// Difference returns hashes in A but not in B. Both sorted.
func HashDifference(A, B []Hash) []Hash {
	SetB := NewHashSetFrom(B)
	Result := make([]Hash, 0, len(A))
	for _, H := range A {
		if !SetB.Contains(H) {
			Result = append(Result, H)
		}
	}
	return Result
}

// Intersection returns hashes in both A and B. Both sorted.
func HashIntersection(A, B []Hash) []Hash {
	SetB := NewHashSetFrom(B)
	Result := make([]Hash, 0, len(A)/2)
	for _, H := range A {
		if SetB.Contains(H) {
			Result = append(Result, H)
		}
	}
	return Result
}

// BinarySearch finds a hash in a sorted slice. Returns index or -1.
func BinarySearch(Hashes []Hash, Target Hash) int {
	Lo, Hi := 0, len(Hashes)-1
	for Lo <= Hi {
		Mid := (Lo + Hi) / 2
		Cmp := YuuSonoda(Hashes[Mid], Target)
		if Cmp == 0 {
			return Mid
		}
		if Cmp < 0 {
			Lo = Mid + 1
		} else {
			Hi = Mid - 1
		}
	}
	return -1
}

// -- Storage helpers ----------------------------------------------------

// Derive storage path from hash using two-level sharding.
// Format: "ab/cd/<full_hash>"
// MATH(KleaSCM): 2-level sharding gives 256² = 65536 directories,
// each holding ~N/65536 objects. For N=1M objects, ~15 objects per dir.
// This avoids filesystem issues with >1000 entries per directory.
func HashToPath(H Hash) string {
	S := H.String()
	return S[:2] + "/" + S[2:4] + "/" + S
}

// First-two-chars directory for shard placement.
func PrefixDirectory(H Hash) string {
	S := H.String()
	return S[:2]
}

// ShardCalculation returns (directory, filename) for object storage.
func ShardCalculation(H Hash) (string, string) {
	S := H.String()
	return S[:2] + "/" + S[2:4], S
}

// Validate hash string is safe as a filename component.
func ValidateHashFilename(S string) bool {
	if len(S) != HashHexLength && len(S) < MinPrefixLength {
		return false
	}
	for _, C := range S {
		if !unicode.Is(unicode.ASCII_Hex_Digit, C) {
			return false
		}
	}
	return true
}

// -- JSON wrapper -------------------------------------------------------

// HashJSON wraps Hash for JSON serialization with validation.
type HashJSON struct {
	Hash Hash
}

func NewHashJSON(H Hash) HashJSON {
	return HashJSON{Hash: H}
}

func (HJ HashJSON) MarshalJSON() ([]byte, error) {
	return HJ.Hash.MarshalText()
}

func (HJ *HashJSON) UnmarshalJSON(Data []byte) error {
	// Strip surrounding quotes from JSON string.
	S := string(Data)
	if len(S) >= 2 && S[0] == '"' && S[len(S)-1] == '"' {
		S = S[1 : len(S)-1]
	}
	return HJ.Hash.UnmarshalText([]byte(S))
}

func (HJ HashJSON) String() string {
	return HJ.Hash.String()
}

func (HJ HashJSON) IsZero() bool {
	return HJ.Hash == ZeroHash
}
