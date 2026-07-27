# Steria Implementation Checklist

---

## 1. HASH — Content addressing core

[x] Hash type with [32]byte fixed-size array
[x] HashFromBytes(b []byte) — SHA-256 single-shot
[x] HashFromReader(r io.Reader) — SHA-256 streaming
[x] HashFromString(s string) — parse hex 64-char string
[x] Hash.String() — encode as lowercase hex
[x] Hash.IsZero() — check all 32 bytes == 0
[x] Hash.Equal(a, b Hash) — constant-time comparison
[x] Hash.Cmp(a, b Hash) — lexicographic ordering for sorting
[x] ZeroHash sentinel constant
[x] Hash.MarshalText() — hex text for JSON/encoding
[x] Hash.UnmarshalText() — hex text parsing
[x] Hash validate on unmarshal — reject non-hex, wrong length
[x] Hash pool — reuse Hash values to reduce allocs
[x] Hash sort — sort.Slice for []Hash ordering
[x] Hash prefix match — match abbreviated hashes
[x] Hash prefix disambiguation — error on ambiguous prefix
[x] Hash prefix min length — configurable minimum abbreviation
[x] Collision check — verify no collision on insert
[x] Double-hash for object addressing — hash of hash
[x] Hash in maps — use as map key
[x] Hash serialization — binary encoding for wire protocol
[x] Hash batch compute — hash multiple buffers
[x] ParseHash() returning (Hash, bool) — error-signaling parse
[x] Strict parser option — reject uppercase, whitespace, 0x prefixes — StrictParseHash
[x] Support 0x prefix parsing — ShizuruFujino
[x] Short hash formatting with configurable length — VivioTakamachi
[x] Base encoding helpers — base32 (RioWesker), base64 (NatsukiKirara)
[x] Validate hash length — reject wrong-length strings — TokakuAzuma
[x] Validate prefix length — reject too-short prefixes — SorawoKamikoshi
[x] Validate prefix characters — reject non-hex prefix chars — TorikoNishina
[x] Validate hash strings before filesystem use — path-safe — RallyVincent
[x] Validate non-zero hash requirement — IsZero guard helper — MayHopkins
[x] Hash.Less(h Hash) — method version for sorting
[x] Hash.Equal(h Hash) — method version for equality
[x] Hash.Compare(h Hash) — three-way method
[x] Byte slice comparison helper — compare raw bytes — EqualBytes
[x] MarshalBinary / UnmarshalBinary — binary wire encoding
[x] Avoid allocation in binary encoding — reuse buffer — FutabaAasu
[x] Fixed-size Hash reader/writer helpers — read (Aer), write (Neviril)
[x] HashFile(path string) Hash — hash file directly from disk — TomokaKase
[x] HashReaders(readers ...io.Reader) Hash — hash multiple readers — NatsukiKuga concatenated
[x] Domain separation: ObjectHash(data) — domain-prefixed hash
[x] Domain separation: TreeHash(data) — tree domain
[x] Domain separation: MetadataHash(data) — metadata domain
[ ] Configurable hash algorithm abstraction — pluggable hash
[x] Efficient prefix lookup for sorted hash lists — binary prefix search via PrefixIndex
[x] Prefix index/cache — accelerate repeated prefix lookups — PrefixIndex
[x] Case handling policy — lowercase normalization — NormalizePrefix
[x] Prefix normalization — strip 0x, lowercase — NormalizePrefix
[ ] Address derivation helper — validate derivations
[x] Reverse mapping support — ObjectID from content hash — Kaguya
[x] ObjectID type separate from content hash — ObjectID struct
[x] HashSet — set of Hash with add/contains/remove/union/intersect
[x] HashMap[V] — map wrapper with Hash keys
[x] Deduplicate hash slices — DeduplicateHashes
[x] Merge hash slices — MergeHashes (two-pointer sorted merge)
[x] Difference/intersection helpers — HashDifference, HashIntersection
[x] Binary-search lookup on sorted hash slices — BinarySearch
[x] JSON object wrapper support — HashJSON type
[x] Text JSON validation — HashJSON.UnmarshalJSON validation
[ ] Binary wire format versioning — forward compat header
[ ] Endianness documentation for wire protocols
[x] HashToPath(h Hash) string — storage path derivation
[x] PrefixDirectory(h Hash) string — first-two-chars dir
[x] ShardCalculation(h Hash) (dir, file) — sharded layout
[x] HashFilename validation — ValidateHashFilename

## 2. ARENA — Memory allocation

[x] Arena type with []byte backing slab
[x] Arena.New(size int) — bump alloc
[x] Arena.Alloc — ZeroBlock on exhaustion (ZII — no panics)
[x] Arena.Reset() — reset bump pointer
[x] Arena.FreeBytes() — remaining capacity
[x] Arena.UsedBytes() — current usage
[x] Arena aligned allocation — 8/16 byte aligned
[x] Arena SIMD alignment — 32/64 byte
[x] Arena string allocation — copy into arena
[x] Arena pool — reuse arenas
[x] Arena grow — new slab, copy
[x] Zero-is-valid — never returns nil
[x] Arena overflow detection
[x] Arena page-aligned mmap for large slabs
[x] Arena mmap release on reset
[x] Arena thread-local for goroutines
[x] Arena pointer tracking (debug)
[x] Arena corruption guard pages
[x] Arena iterator — walk allocations

## 3. BLOB — Raw content storage

[x] Blob type wrapping []byte — KuyuMashima
[x] Blob.Hash() — SulettaMercury
[x] Blob.Size() — YuuKoito
[x] Blob.Type() — returns "blob"
[x] Blob.Encode() — MiriamHildegardvonGropius
[x] Blob.Decode() — RiriHitotsuyanagi
[x] Blob compression — gzip before store — Compress method
[x] Blob decompression — gzip on read — Decompress method
[x] Blob compression level — configurable
[x] Blob compression threshold — skip small
[x] Binary blob detection heuristic — MioChibana
[x] Text blob detection — UTF-8 — UshioKazama
[x] Blob dedup — skip if hash exists
[x] Large blob chunking — split into chunks — HarukaTenou,MichiruKaioh
[x] Chunk size — configurable (default 8KB min)
[x] Content-defined chunk boundaries (cyclic polynomial CDC) — HarukaTenou
[x] Chunk index — hash to chunk list — MichiruKaioh
[x] Chunk dedup across blobs — KirikaAkatsuki
[x] Blob streaming read — ShirabeTsukuyomi
[x] Blob streaming write — ChrisYukine
[x] Blob ref counting for GC — TsubasaKazanari,KanadeAmou
[x] Blob pinning — prevent GC — Pin,Unpin
[x] Empty blob handling — MariaCadenzavnaEve
[x] Blob MIME type detection — MasakiAkemiya
[x] Blob language detection — TomoeHachisuka
[x] Blob line ending detection — ReiHino
[x] Blob encoding detection — MinakoAino
[ ] Blob lifecycle — Clone, SetData, Clear, Reload, LazyLoad
[ ] Immutable blob enforcement — reject modification after hash
[ ] Validate blob hash matches data — integrity check
[ ] Validate blob type — reject unknown types
[ ] Validate MIME value — reject malformed MIME
[ ] Validate encoding value — reject unsupported encodings
[ ] Validate chunk index — consistent boundary list
[ ] Validate compressed blob metadata — algorithm + sizes match
[ ] Detect invalid/corrupted blob serialization — checksum
[ ] Deserialize with metadata preservation — round-trip
[ ] Serialize metadata separately — header/data split
[ ] Versioned blob format — forward compat
[ ] Compression flag in serialized format — enable/disable bit
[ ] Chunk metadata serialization — serialize ChunkIndex
[ ] Streaming serialization/deserialization — large blob pipelining
[ ] Store compression metadata: algorithm, original size, compressed size
[ ] Detect compressed blobs automatically — magic bytes
[ ] Support alternate compression: zstd
[ ] Support alternate compression: lz4
[ ] Compression ratio statistics — before/after bytes
[ ] Streaming compression/decompression — chunked pipeline
[ ] Hash compressed vs uncompressed content distinction
[ ] Content hash verification after decompression
[ ] BlobID separate from storage hash — addressing indirection
[ ] Hash reader streams directly — io.Reader → Hash
[ ] Hash file from disk — FileHash(path)
[ ] Better MIME detection — content-type sniffing
[ ] Better language detection — linguist-style heuristics
[ ] Shebang detection — #!/ interpreter detection
[ ] Binary format detection: PNG — magic bytes
[ ] Binary format detection: JPEG — magic bytes
[ ] Binary format detection: ELF — magic bytes
[ ] Binary format detection: ZIP — magic bytes
[ ] Binary format detection: PDF — magic bytes
[ ] Charset detection — detect encoding from content
[ ] BOM handling — strip/preserve Byte Order Mark
[ ] Mixed line ending detection — CR+LF vs LF vs CR
[ ] Store chunk boundaries in ChunkIndex — persistent index
[ ] Reconstruct blob from chunks — chunk-to-blob assembly
[ ] Verify chunk hashes — integrity per chunk
[ ] Serialize chunk indexes — wire format for chunk list
[ ] Chunk size validation — min/max configurable
[ ] Different chunking algorithms — gear hash, buzhash
[ ] Rolling hash state persistence — resume CDC
[ ] Store chunks independently — chunk objects in store
[ ] Retrieve chunks by hash — chunk lookup
[ ] Missing chunk detection — detect gaps
[ ] Chunk garbage collection — orphan chunk cleanup
[ ] Chunk reference counting — track chunk users
[ ] Chunk reuse statistics — dedup ratio tracking
[ ] Blob diff — difference between two blobs
[ ] Text line diff — line-by-line comparison
[ ] Binary diff — byte-level diff
[ ] Patch generation — create patch from diff
[ ] Patch application — apply patch to blob
[ ] Similarity calculation — hash-based similarity
[ ] Text search indexing — full-text index
[ ] Extracted metadata indexing — MIME/language index
[ ] Content fingerprinting — minhash/simhash
[ ] Duplicate blob detection — exact + near-duplicate
[ ] GC scan references — find all blob references
[ ] GC find unreferenced blobs — mark-sweep
[ ] GC delete unused blobs — sweep phase
[ ] GC respect pinned blobs — skip pinned
[ ] GC recalculate references — ref count audit
[ ] GC verify ref counts — consistency check
[ ] Streaming blob creation — write-as-you-go
[ ] Avoid copying data unnecessarily — zero-copy paths
[ ] Memory limit handling — cap large blob memory
[ ] Large blob support — multi-GB
[ ] Blob cache — in-memory LRU
[ ] Chunk cache — chunk-level cache
[ ] Decompression bomb protection — size limit
[ ] Maximum blob size limit — configurable
[ ] Maximum decompressed size limit — safety limit
[ ] Safe MIME handling — reject dangerous MIME types
[ ] Safe language detection limits — timeout/cap
[ ] IsEmpty() — blob has no data
[ ] IsCompressed() — blob uses compression
[ ] Size() — data size
[ ] CompressedSize() — stored size
[ ] ChunkCount() — number of chunks
[ ] HasChunks() — blob is chunked
[ ] ContentHash() — hash of content
[ ] Metadata() — all metadata access
[ ] Write blob to object store — store integration
[ ] Read blob from object store — load integration
[ ] Store chunked blobs — multi-object chunk store
[ ] Restore chunked blobs — chunk-to-blob restore
[ ] Verify stored blob — integrity after store

## 4. TREE — Directory representation

[x] Tree type with sorted entry list — HomuraAkemi, FateTestarossa
[x] TreeEntry with Mode, Path, Hash — TreeEntry struct
[x] TreeEntry mode — file, executable, symlink, dir — TreeModeFile/Exec/Symlink/Dir/Submodule
[x] TreeEntry mode bits — Unix permissions — SanyaLitvyak, PerrineNoel
[x] Tree canonical sort order — sort() + entryLess
[x] Tree sort by path byte-wise — entryLess with implicit '/' on dirs
[x] Tree.Hash() — recursive hash — YoshikaMiyafuji
[x] Tree.Encode() — flat serialization — EricaHartmann
[x] Tree.Decode() — parse — MioSakamoto
[x] Tree walk — depth-first traversal — MinnaWilcke
[x] Tree lookup — entry by path — EilaIlmatarJuutilainen (binary search)
[x] Tree diff — compare two trees — CharlotteYeager
[x] Tree merge — combine + conflict detection — HonamiIchinose
[x] Tree filter — exclude paths — LynneBowman, AlisaOmela
[x] Tree from filesystem — walk directory — NipaJeanne
[x] Tree to filesystem — write to disk — WilckeJager
[x] Tree entry mode for submodules — TreeModeSubmodule constant
[x] Tree path encoding — UTF-8 — MiyabiUsami
[x] Tree path validation — reject null/absolute — SylvetteChanel
[x] Tree size estimation — NozomiHitomi
[x] Tree entry count — TiltyClaret
[x] Sparse tree — partial — EriKawamura
[x] Tree cache — FIFO — TreeCache, MiyakoMiyamura, RikaMatsumoto
[x] Tree delta — compressed diff — SuzuneHorikita
[x] Tree delta apply — KeiKaruizawa
[x] Empty tree — zero entries — ZeroTree, IsZero
[x] Tree merge recursive — nested — HonamiIchinose (auto-resolves non-conflicting)
[x] Tree conflict markers — HonamiIchinose (returns Conflicts list)
[x] Subtree extraction — MaeMatsuo
[x] Tree graft — replace subtree — UrsulaHartmann
[ ] Tree validation — Validate() full tree integrity
[ ] Duplicate path detection — reject duplicate entries
[ ] Parent/child path conflict detection — reject file/dir conflicts
[ ] Mode validation — reject invalid mode values
[ ] Path normalization/validation — normalize separators
[ ] Binary search comparator matching tree sort order — FIX: order bug
[ ] Verify delta before applying — pre-apply integrity
[ ] FindPath(path) (TreeEntry, bool) — explicit find API
[ ] Has(path) — path existence check
[ ] Insert(path, mode, hash) — insert entry
[ ] Remove(path) — remove entry
[ ] Replace(path, mode, hash) — replace entry
[ ] Rename(old, new path) — rename entry
[ ] Clear() — remove all entries
[ ] Equal(other Tree) — structural equality
[ ] Children(dir) — direct children of directory
[ ] WalkPrefix(prefix) — entries under prefix
[ ] Files() — only file entries
[ ] Directories() — only directory entries
[ ] Recursive walk — full depth-first
[ ] Compute file hashes during import — hash on fs walk
[ ] Restore file contents — write file from hash
[ ] Restore permissions — chmod on restore
[ ] Restore symlinks — symlink creation
[ ] Restore timestamps (optional) — utimes
[ ] Ignore/exclude support — skip ignored paths
[ ] Better three-way merge — recursive resolve
[ ] Add/add conflict handling — both add same path
[ ] Delete/delete conflict handling — both delete
[ ] File/directory conflict handling — type change
[ ] Mode conflict handling — permission change
[ ] Rich conflict information — conflict metadata
[ ] Multiple delta hunks — split large delta
[ ] Delta verification — pre-apply checksum
[ ] Tree builder — incremental construction
[ ] Prefix binary search — O(log n) path lookup
[ ] LRU cache (instead of FIFO) — configurable eviction
[ ] Bulk insert/remove operations — batch mutations
[ ] EntryCount — total entries
[ ] FileCount — file entries
[ ] DirectoryCount — directory entries
[ ] MaximumDepth — tree depth
[ ] TreeSize — serialized size estimate
[ ] IsFile(entry) — check if file mode
[ ] IsDir(entry) — check if dir mode
[ ] IsExec(entry) — check if executable
[ ] IsSymlink(entry) — check if symlink
[ ] IsSubmodule(entry) — check if submodule
[ ] Parent directory helpers — parent path from entry
[ ] Subtree replace helpers — replace entry subtree
[ ] Bulk rename/move — batch rename entries
[ ] MarshalBinary / UnmarshalBinary — binary serialization
[ ] Serialization validation — round-trip verify
[ ] Unit tests — comprehensive test coverage
[ ] Fuzz tests — random tree operations
[ ] Property tests — invariant testing
[ ] Benchmarks — performance regression tests
[ ] FIXED: Binary search uses different ordering than entryLess — lookups fail for correctly sorted trees

## 5. COMMIT — Snapshot with metadata

[ ] Commit type — tree, parent(s), author, committer, message, timestamp
[ ] Single parent — normal commit
[ ] Multiple parents — merge commit
[ ] Zero parents — root commit
[ ] Commit.Hash() — content-addressable
[ ] Commit.Encode() — serialized
[ ] Commit.Decode() — parse
[ ] Commit tree link
[ ] Commit author — identity + timestamp
[ ] Commit committer — distinct from author
[ ] Commit message — subject + body
[ ] Commit message encoding — UTF-8
[ ] Commit message wrap — 72 chars
[ ] Commit timestamp — Unix + timezone
[ ] Commit timestamp validation
[ ] Commit signature — inline GPG/SSH
[ ] Commit signature verification
[ ] Commit signature format — RFC 4880
[ ] Commit parent ordering
[ ] Commit generation number
[ ] Commit topological sort
[ ] Commit graph — DAG
[ ] Commit graph file storage
[ ] Commit graph generation
[ ] Commit graph Bloom filter
[ ] Commit message template
[ ] Commit message editor — $EDITOR
[ ] Commit message cleanup
[ ] Commit amend — replace HEAD
[ ] Commit fixup — autosquash marker
[ ] Commit squash — combine
[ ] Commit reword — change message
[ ] Commit sign — auto-sign
[ ] Commit verify — verify range
[ ] Commit grep — search messages
[ ] Commit range — .. and ...
[ ] Commit ancestors — parent walk
[ ] Commit descendants — child walk
[ ] Commit message trailers
[ ] Commit co-authored-by
[ ] Commit message length limit
[ ] Commit empty message reject

## 6. SUPERPOSITION INDEX — Core data structure

[ ] Index type — file path to version list
[ ] Index entry — path, mode, hash list, current version
[ ] Index version — hash, timestamp, identity, message
[ ] Index add file — insert/update
[ ] Index remove file — delete entry
[ ] Index lookup — by path
[ ] Index contains — path check
[ ] Index count — tracked files
[ ] Index iterator — walk entries
[ ] Index snapshot — immutable copy
[ ] Index version append
[ ] Index version current — active
[ ] Index version rollback
[ ] Index superposition state — multi-version
[ ] Index collapse — resolve to single
[ ] Index save — serialize to .steria/index
[ ] Index load — parse from file
[ ] Index format — binary encoding
[ ] Index format version — forward compat
[ ] Index migrate — upgrade format
[ ] Index lock — concurrent write prevention
[ ] Index merge — combine indices
[ ] Index diff — differences
[ ] Index conflict — superposition conflict
[ ] Index conflict resolve
[ ] Index checksum — integrity
[ ] Index corruption detection
[ ] Index rebuild from objects
[ ] Index size limit
[ ] Index version limit per file
[ ] Index compaction
[ ] Index filter — include/exclude
[ ] Index wildcard — glob
[ ] Index status — added/modified/removed/conflict
[ ] Index sync status — ahead/behind
[ ] Index parent — link to previous
[ ] Index branch — which branch
[ ] Index add file with hash — AddFile(path, hash, mode)
[ ] Index remove file — RemoveFile(path)
[ ] Index rename/move file — RenameFile(old, new)
[ ] Index copy file history — CopyHistory(src, dst)
[ ] Index replace file entry — ReplaceFile(path, entry)
[ ] Index clear — remove all entries
[ ] Index merge — Merge(other) combine indices
[ ] Index Has(path) — quick path existence
[ ] Index Find(path) — get entry by path
[ ] Index FindHash(hash) — find entries by content hash
[ ] Index Latest(path) — most recent version entry
[ ] Index LatestHash(path) — hash of latest version
[ ] Index LatestVersion(path) — version number of latest
[ ] Index Version(path, idx) — get specific version
[ ] Index VersionByHash(hash) — get version by hash
[ ] Index AppendVersion(path, version) — append to history
[ ] Index RemoveVersion(path, idx) — remove from history
[ ] Index ReplaceVersion(path, idx, version) — replace version
[ ] Index Rollback(path, idx) — roll back to previous
[ ] Index RestoreVersion(path, idx) — restore specific version
[ ] Index TruncateHistory(path, N) — keep N most recent
[ ] Index SquashHistory(path) — collapse to single
[ ] Index QueryByHash(hash) — search versions by hash
[ ] Index QueryByTimestamp(ts) — search by timestamp
[ ] Index QueryByIdentity(id) — search by author
[ ] Index QueryByMessage(msg) — search by commit message
[ ] Index QueryByPath(pattern) — search by path glob
[ ] Index SetHEAD(ref) — set HEAD reference
[ ] Index GetHEAD() — get HEAD reference
[ ] Index ClearHEAD() — clear HEAD
[ ] Index ValidateHEAD() — verify HEAD exists
[ ] Index HEADExists() — check HEAD presence
[ ] Index Validate() — full index integrity
[ ] Index ValidateFilePaths() — reject invalid paths
[ ] Index ValidateHashes() — reject corrupt hashes
[ ] Index ValidateTimestamps() — reject impossible timestamps
[ ] Index ValidateIdentities() — reject missing identity
[ ] Index DetectDuplicateEntries() — find duplicate paths
[ ] Index DetectDuplicateVersions() — find duplicate version hashes
[ ] Index DetectOrphanedVersions() — versions with no owner
[ ] Index DetectInvalidHEAD() — HEAD points nowhere
[ ] Index DetectMissingObjectRefs() — hashes not in object store
[ ] Index CompareIndexes(a, b) — diff two indices
[ ] Index CompareFileHistories(path, a, b) — diff file histories
[ ] Index CompareVersions(a, b) — diff two versions
[ ] Index ListAdded() — files added relative to base
[ ] Index ListRemoved() — files removed relative to base
[ ] Index ListModified() — files changed relative to base
[ ] Index VerifyIntegrity() — checksum verification
[ ] Index RebuildFromObjects() — reconstruct from store
[ ] Index Compact() — reduce storage size
[ ] Index Optimize() — performance optimization
[ ] Index Repair() — fix corruption
[ ] Index AtomicWrite() — temp + rename
[ ] Index CrashSafeUpdate() — journaled writes
[ ] Index Backup() — snapshot before overwrite
[ ] Index Recover() — recovery after interrupted write
[ ] Fast path lookup — hash-indexed path map
[ ] Fast hash lookup — path-indexed hash map
[ ] Fast version lookup — version-indexed
[ ] Lazy history loading — load on demand
[ ] Incremental writes — append-only log
[ ] Optional in-memory cache — performance
[ ] Concurrent readers — RLock
[ ] Safe concurrent writers — exclusive lock
[ ] Repository locking — file lock
[ ] Duplicate write protection — compare-and-swap
[ ] FileCount — tracked files
[ ] VersionCount — total versions
[ ] UniqueIdentities — unique committers
[ ] OldestVersion — earliest timestamp
[ ] NewestVersion — latest timestamp
[ ] LargestHistory — most versions for a file
[ ] AverageHistoryDepth — mean versions per file
[ ] RepositoryAge — time span
[ ] IterateFiles() — iterate all entries
[ ] IterateVersions() — iterate all versions
[ ] IterateChronological() — ordered by time
[ ] IterateByPrefix(prefix) — path-prefix iteration
[ ] ExportIndex() — portable format
[ ] ImportIndex(data) — load portable format
[ ] ExportHistory(path) — export file history
[ ] ImportHistory(path, data) — import file history
[ ] Serialize() — binary encoding
[ ] Deserialize(data) — binary decoding
[ ] SearchByPath(pattern) — path search
[ ] SearchByMessage(text) — message full-text
[ ] SearchByIdentity(name) — author/committer
[ ] SearchByHash(hash) — content hash search
[ ] SearchByTimeRange(start, end) — time-range search
[ ] PruneUnreachableVersions() — sweep orphaned data
[ ] GCUnusedObjects() — reclaim object store space
[ ] RemoveEmptyHistories() — delete empty entries
[ ] DeduplicateVersions() — merge identical versions
[ ] VerifyObjectReferences() — cross-check with store
[ ] ListTrackedFiles() — all tracked paths
[ ] ListTrackedHashes() — all referenced hashes
[ ] ListIdentities() — all unique identities
[ ] RepositorySummary() — stats report
[ ] CloneIndex() — deep copy index
[ ] Equality() — full structural equality

## 7. SUPERPOSITION MODEL

[ ] Superposition fork — parallel version
[ ] Superposition collapse — select one
[ ] Superposition merge — combine lines
[ ] Superposition observe — keep version
[ ] Superposition decohere — discard
[ ] Superposition state — in index entry
[ ] Superposition visualization — version tree
[ ] Superposition propagation — cascade
[ ] Superposition entanglement — linked
[ ] Superposition measurement — deterministic
[ ] Superposition interference — detect conflict
[ ] Superposition basis — orthogonal axes
[ ] Superposition collapse strategy — random/weighted/manual
[ ] Superposition collapse dry-run
[ ] Superposition collapse undo
[ ] Superposition history — log
[ ] Superposition log — timeline
[ ] Superposition diff — compare states
[ ] Superposition rebase
[ ] Superposition cherry-pick
[ ] Superposition branch create
[ ] Superposition branch delete
[ ] Superposition branch list
[ ] Superposition branch switch
[ ] Superposition branch rename
[ ] Superposition branch merge
[ ] Superposition branch rebase
[ ] Superposition branch tracking — remote
[ ] Superposition divergence detection
[ ] Superposition convergence detection
[ ] Superposition complexity metric
[ ] Superposition probability — weight
[ ] Superposition amplitude — confidence
[ ] Multiple superposition branches
[ ] Superposition branch global state

## 8. OBJECT STORE

[ ] Object store interface — Read, Write, Exists, Delete, Iterate
[ ] Object type discrimination
[ ] Object write — hash, store by hash
[ ] Object read — retrieve by hash
[ ] Object exists — check presence
[ ] Object delete — remove
[ ] Object iterate — enumerate
[ ] Object size — stored size
[ ] Object compression — zstd
[ ] Object decompression
[ ] Object cache — LRU
[ ] Object cache size config
[ ] Object cache TTL
[ ] Object ref count
[ ] Object reachability from refs
[ ] Object pack — batch into pack
[ ] Object unpack — extract
[ ] Object verify — hash integrity
[ ] Object corrupt — mark
[ ] Object repair — alternate source
[ ] Object quarantine — hold until commit
[ ] Object quarantine promote
[ ] Object quarantine discard
[ ] Object store stats — count, size
[ ] Object store path — hash prefix dirs
[ ] Object store migration
[ ] Object store backup — incremental
[ ] Object store restore
[ ] Object store FS — filesystem
[ ] Object store memory — in-memory
[ ] Object store S3 — AWS backend
[ ] Object store plugin interface
[ ] Object write atomic — temp + rename
[ ] Object write retry
[ ] Object write fsync
[ ] DeleteObject(hash) — remove object from store
[ ] ListObjects() — enumerate all objects
[ ] ObjectSize(hash) — stored size
[ ] VerifyObjectHash(hash) — content matches hash
[ ] DetectCorruptedObject — checksum validation
[ ] DetectTruncatedObject — size mismatch
[ ] Atomic writes — write temp file then rename
[ ] Crash-safe writes — avoid partial writes
[ ] PruneUnreachableObjects — GC sweep
[ ] GarbageCollect — full GC cycle
[ ] ObjectStats — count, bytes used, type breakdown
[ ] StorageUsageReport — per-type storage
[ ] Avoid repeated MkdirAll — cache dir existence
[ ] Buffered I/O for large objects — chunked read/write
[ ] Optional object cache — in-memory LRU
[ ] Batch read operations — multi-object read
[ ] Batch write operations — multi-object write
[ ] Safe concurrent writes — mutex/mmap sync
[ ] File locking — cross-process lock
[ ] Avoid duplicate write races — compare-and-store
[ ] Validate object path — reject path traversal
[ ] Reject invalid hashes before filesystem access — early validation
[ ] Read into existing buffer — reuse allocation
[ ] Streaming reader (io.Reader) — sequential read
[ ] Streaming writer (io.Writer) — sequential write
[ ] Memory-mapped reads — mmap for large objects (optional)
[ ] Compress on write — transparent compression
[ ] Decompress on read — transparent decompression
[ ] Configurable compression level — speed/size tradeoff
[ ] InitializeObjectDirectory — create storage layout
[ ] VerifyObjectDatabase (fsck) — full integrity scan
[ ] RepairMissingDirectories — recreate paths
[ ] CountObjects — total stored count

## 9. PACK FILE

[ ] Pack header — signature, version, count
[ ] Pack object entry — type, size, compressed data
[ ] Pack object types — blob, tree, commit, tag
[ ] Pack object offsets
[ ] Pack CRC — per-object checksum
[ ] Pack SHA — file checksum
[ ] Pack index (.idx) — hash-to-offset
[ ] Pack index fanout — 256 prefix entries
[ ] Pack index lookup — binary search
[ ] Pack index CRC
[ ] Pack index version
[ ] Pack multi-index
[ ] Pack thin — without bases
[ ] Pack thin fix — add bases on receive
[ ] Pack delta — store differences
[ ] Pack delta instruction — copy/insert
[ ] Pack delta apply — reconstruct
[ ] Pack delta chain — N-level limit
[ ] Pack delta compression — find base
[ ] Pack delta selection — greedy/windowed
[ ] Pack delta window — candidates
[ ] Pack delta depth — max chain
[ ] Pack progress — creation progress
[ ] Pack memory — budget
[ ] Pack thread — parallel compression
[ ] Pack verify — all objects + deltas
[ ] Pack unpack — extract to loose
[ ] Pack repack — consolidate
[ ] Pack full — no deltas
[ ] Pack append — add objects
[ ] Pack rewrite — optimize
[ ] Pack prune — remove objects
[ ] Pack list — enumerate
[ ] Pack size
[ ] Pack object count
[ ] Pack duplicate detection
[ ] Pack sort — locality ordering
[ ] Pack fanout — multi-pack dispatch
[ ] Pack data mmap
[ ] Pack data sliding window
[ ] Pack delta base cache
[ ] Pack delta base cache size
[ ] Pack keep — prevent repack
[ ] Pack keep file — .keep marker

## 10. REFERENCE — Named pointers

[ ] Ref type — name to hash
[ ] Ref resolve — follow symbolic
[ ] Ref name — full with namespace
[ ] Ref short name — strip prefix
[ ] Ref parse — validate format
[ ] Ref validate — reject invalid chars
[ ] Ref namespace — heads, tags, remotes
[ ] Ref loose — file per ref
[ ] Ref packed — single file
[ ] Ref packed parse
[ ] Ref packed update — unpack on modify
[ ] Ref symbolic — pointer to ref
[ ] Ref symbolic resolve — follow chain
[ ] Ref create — write
[ ] Ref delete — remove
[ ] Ref update — change hash
[ ] Ref rename
[ ] Ref list — enumerate by namespace
[ ] Ref match — glob pattern
[ ] Ref transaction — atomic batch
[ ] Ref lock — prevent concurrent
[ ] Ref log (reflog) — history
[ ] Reflog entry — old/new hash, committer, message
[ ] Reflog read — parse file
[ ] Reflog write — append
[ ] Reflog expire — remove old
[ ] Reflog gc — cleanup
[ ] Reflog limit — max entries
[ ] HEAD ref — current branch
[ ] HEAD resolve — current commit
[ ] ORIG_HEAD — previous HEAD
[ ] FETCH_HEAD — last fetch
[ ] MERGE_HEAD — merge sources
[ ] CHERRY_PICK_HEAD
[ ] BISECT_HEAD
[ ] Refspec — remote mapping
[ ] Refspec parse — +src:dst
[ ] Refspec match — source pattern
[ ] Refspec expand — generate destination
[ ] Refspec force — non-ff allowed
[ ] Note ref — refs/notes/

## 11. STERIA DIRECTORY — .steria layout

[ ] .steria/ on watch/init
[ ] .steria/config — project config
[ ] .steria/index — superposition index
[ ] .steria/objects/ — loose objects
[ ] .steria/packs/ — pack files
[ ] .steria/refs/heads/ — branches
[ ] .steria/refs/tags/ — tags
[ ] .steria/refs/remotes/ — remote tracking
[ ] .steria/refs/super/ — superposition refs
[ ] .steria/HEAD — current branch
[ ] .steria/reflog/ — reference logs
[ ] .steria/hooks/ — hook scripts
[ ] .steria/hooks/pre-done
[ ] .steria/hooks/post-done
[ ] .steria/hooks/pre-choose
[ ] .steria/hooks/post-choose
[ ] .steria/hooks/pre-merge
[ ] .steria/hooks/post-merge
[ ] .steria/hooks/pre-push
[ ] .steria/hooks/post-push
[ ] .steria/hooks/pre-pull
[ ] .steria/hooks/post-pull
[ ] .steria/info/exclude — local ignore
[ ] .steria/info/attributes
[ ] .steria/info/refs — dumb protocol
[ ] .steria/description
[ ] .steria/superposition/ — state
[ ] .steria/lock — repo lock
[ ] .steria/index.lock
[ ] .steria/prepare-commit-msg
[ ] .steria/COMMIT_EDITMSG
[ ] .steria/MERGE_HEAD
[ ] .steria/MERGE_MSG
[ ] .steria/logs/HEAD — HEAD reflog

## 12. WATCH — Initialize tracking

[ ] Watch command — init .steria
[ ] Check existing .steria — error
[ ] Create .steria directory structure
[ ] Write initial config with defaults
[ ] Create empty index
[ ] Set HEAD to refs/heads/main
[ ] Create main branch ref
[ ] Create initial empty commit
[ ] Set HEAD to initial commit
[ ] Check user identity configured
[ ] Default branch name configurable
[ ] Bare mode — no working tree
[ ] Template directory support
[ ] Shared permissions mode
[ ] Quiet mode
[ ] Verbose mode
[ ] Output .steria path
[ ] Detect existing files
[ ] Auto-add existing files option
[ ] Set initial .steriaignore
[ ] Configure remote origin
[ ] Set repo description
[ ] Install default hooks
[ ] Dry-run mode
[ ] Force reinitialize
[ ] External .steria location
[ ] Worktree support
[ ] Directory permissions
[ ] Race detection — concurrent init
[ ] Rollback on failure
[ ] Log init event

## 13. DONE — Save snapshot

[ ] Done command — create snapshot
[ ] Find working tree files from index
[ ] Filter by path — specific files
[ ] Filter by glob pattern
[ ] Include untracked files
[ ] Scan filesystem for changes
[ ] Hash all files
[ ] Compare with index — detect changes
[ ] Dedup unchanged — skip same hash
[ ] Detect new files — add to index
[ ] Detect deleted files — remove
[ ] Detect renames — similarity
[ ] Detect copies — hash match
[ ] Detect mode changes — permissions
[ ] Build tree from index
[ ] Write tree object to store
[ ] Create commit object
[ ] Set parent from HEAD
[ ] Update HEAD to new commit
[ ] Update current branch ref
[ ] Write reflog entry
[ ] Commit message from -m
[ ] Commit message from file -F
[ ] Commit message from editor
[ ] Commit message template
[ ] Commit message cleanup
[ ] Commit with diff in editor
[ ] Commit with status in editor
[ ] Allow empty commit
[ ] Amend previous commit
[ ] Sign commit with key
[ ] Author override
[ ] Date override
[ ] Committer override
[ ] Auto-add tracked modified -a
[ ] Include paths -i
[ ] Exclude paths
[ ] Renormalize attributes
[ ] No-verify — skip hooks
[ ] Dry-run
[ ] Status display
[ ] Verbose diff output
[ ] Quiet mode
[ ] Progress for large commits
[ ] Auto-done on filesystem events
[ ] Debounce auto-done
[ ] Cleanup temp files
[ ] Atomic — all-or-nothing
[ ] Rollback on failure
[ ] Lock — prevent concurrent
[ ] Pre-done hook
[ ] Post-done hook
[ ] Hook failure abort
[ ] Conventional commit scopes
[ ] Commit trailers — Signed-off-by
[ ] Co-authored-by support

## 14. CHOOSE — Collapse superposition

[ ] Choose command — resolve superposition
[ ] Choose list — files in superposition
[ ] Choose status — summary
[ ] Choose file — collapse specific
[ ] Choose pattern — glob match
[ ] Choose all — collapse all
[ ] Choose version by index
[ ] Choose version by hash
[ ] Choose version by date
[ ] Choose version by message
[ ] Choose interactive — per-file
[ ] Choose diff — show version diff
[ ] Choose merge — merge versions
[ ] Choose ours — local version
[ ] Choose theirs — remote version
[ ] Choose auto — heuristic
[ ] Choose weighted — probability
[ ] Choose dry-run
[ ] Choose undo
[ ] Choose log — history
[ ] Choose branch — collapse branch
[ ] Choose abort
[ ] Choose continue
[ ] Choose skip — leave unresolved
[ ] Choose mark resolved
[ ] Choose conflict detection
[ ] Choose external merge tool
[ ] Choose update working tree
[ ] Choose update index
[ ] Choose create commit
[ ] Choose no-commit
[ ] Choose strategy — configurable
[ ] Choose strategy — patience
[ ] Choose strategy — histogram
[ ] Choose strategy — minimal
[ ] Choose rename detection
[ ] Choose binary handling
[ ] Choose symlink handling
[ ] Choose submodule handling
[ ] Choose whitespace handling
[ ] Choose conflict style — diff3/zdiff3
[ ] Choose pre-hook
[ ] Choose post-hook
[ ] Choose progress
[ ] Choose conflict markers
[ ] Choose backup conflicting

[ ] Choose no-edit
[ ] Choose cleanup -- message cleanup
[ ] Choose rerere-autoupdate
[ ] Choose verbosity levels
[ ] Choose path exclusion
[ ] Choose path inclusion by pattern
[ ] Choose annotated markers -- file:line annotations
[ ] Choose merge tool trust exit code
[ ] Choose merge tool keep backup
[ ] Choose merge tool keep temporaries
[ ] Choose merge tool gui preference
[ ] Choose merge tool prompt
[ ] Choose merge tool list available
[ ] Choose merge tool configuration per tool
[ ] Choose diff algorithm selection per file
[ ] Choose strategy per file via attributes

## 15. INIT — Remote/project setup

[ ] Init command
[ ] Init local — create .steria
[ ] Init bare — no working tree
[ ] Init remote — create on server
[ ] Init remote URL
[ ] Init remote via SSH
[ ] Init remote via HTTP
[ ] Init with template
[ ] Init shared permissions
[ ] Init with description
[ ] Init with hooks
[ ] Init with config values
[ ] Init quiet
[ ] Init verbose
[ ] Init existing check
[ ] Init force
[ ] Init dry-run
[ ] Init remote auth config
[ ] Init remote token
[ ] Init remote SSH key
[ ] Init remote config
[ ] Init object format config
[ ] Init ref format config
[ ] Init default branch name
[ ] Init worktree add
[ ] Init separate git dir
[ ] Init disk space check
[ ] Init permissions check
[ ] Init error messages
[ ] Init rollback
[ ] Init log event

## 16. CLONE — Copy remote

[ ] Clone command
[ ] Clone full — all objects
[ ] Clone shallow — limited depth
[ ] Clone depth N
[ ] Clone since date
[ ] Clone single branch
[ ] Clone bare
[ ] Clone mirror — exact replica
[ ] Clone with tags
[ ] Clone no-tags
[ ] Clone recursive — submodules
[ ] Clone submodule depth
[ ] Clone to directory
[ ] Clone to existing error
[ ] Clone to existing force
[ ] Clone create .steria
[ ] Clone fetch remote refs
[ ] Clone negotiate objects
[ ] Clone receive pack
[ ] Clone unpack objects
[ ] Clone checkout default branch
[ ] Clone checkout specific branch
[ ] Clone no-checkout
[ ] Clone progress bar
[ ] Clone quiet
[ ] Clone verbose
[ ] Clone filter — partial
[ ] Clone filter blob:none
[ ] Clone filter tree:0
[ ] Clone filter blob:limit=N
[ ] Clone filter spec — sparse:oid=
[ ] Clone remote named origin
[ ] Clone remote URL store
[ ] Clone remote auth
[ ] Clone resume interrupted
[ ] Clone retry on failure
[ ] Clone timeout
[ ] Clone verify objects
[ ] Clone reference repo
[ ] Clone server options
[ ] Clone config values
[ ] Clone template
[ ] Clone hooks setup
[ ] Clone shallow since
[ ] Clone shallow exclude refs
[ ] Clone unshallow
[ ] Clone dissociate from reference
[ ] Clone atomic

## 17. SERVE — Remote daemon

[ ] Serve command — HTTP daemon
[ ] Bind address config
[ ] Port config
[ ] Repository root path
[ ] Serve bare repos
[ ] Auth middleware
[ ] Auth basic
[ ] Auth token bearer
[ ] Auth none
[ ] TLS support
[ ] TLS cert path
[ ] TLS key path
[ ] TLS auto Let's Encrypt
[ ] Allow push
[ ] Allow pull
[ ] Allow create repo
[ ] Access log — Apache format
[ ] Error log
[ ] Log format config
[ ] Log file output
[ ] Log stdout
[ ] Rate limit per IP
[ ] Rate limit burst
[ ] Max body size
[ ] Max objects per request
[ ] Read timeout
[ ] Write timeout
[ ] Idle timeout
[ ] Max header size
[ ] Max connections
[ ] CORS headers
[ ] CORS origins
[ ] CORS methods
[ ] CORS headers allowed
[ ] Gzip compression
[ ] Compression level
[ ] Advertise repos — GET /
[ ] Info/refs — smart protocol
[ ] Receive pack — push
[ ] Upload pack — pull
[ ] Stateless mode
[ ] Repo create — POST /repos
[ ] Repo list — GET /repos
[ ] Repo delete — DELETE /repos
[ ] Repo info — GET /repos/:name
[ ] Health check — GET /health
[ ] Metrics — GET /metrics
[ ] Graceful shutdown — SIGTERM
[ ] Daemonize — background
[ ] PID file write
[ ] Drop privileges user
[ ] Drop privileges group
[ ] Chroot jail
[ ] Private — deny anonymous
[ ] Public — allow anonymous read
[ ] Repo visibility check
[ ] Request body close
[ ] Request ID middleware
[ ] Panic recover middleware
[ ] Content-Type enforcement
[ ] Response status codes
[ ] Response JSON encoder
[ ] Response error body format
[ ] Security headers
[ ] Query parameter parsing
[ ] Repo-level permissions

## 18. CONFIG — Configuration system

[ ] Config command — get/set/list/unset
[ ] Config get by key
[ ] Config set by key
[ ] Config unset by key
[ ] Config list all
[ ] Config list regex filter
[ ] Config global — ~/.steria/config
[ ] Config local — .steria/config
[ ] Config system — /etc/steria/config
[ ] Config file format JSON
[ ] Config file permissions 0600
[ ] Config file create on first set
[ ] Config file backup before modify
[ ] Config layers — global < local < system
[ ] Config env override
[ ] Config type — string, bool, int, path
[ ] Config bool parsing — true/false/1/0
[ ] Config int parsing
[ ] Config path normalization
[ ] Config key format — section.name
[ ] Config key validation
[ ] Config value validation
[ ] Config defaults — built-in
[ ] Config user.name
[ ] Config user.email
[ ] Config user.signingKey
[ ] Config core.editor
[ ] Config core.pager
[ ] Config core.autocrlf
[ ] Config core.safecrlf
[ ] Config core.ignorecase
[ ] Config core.protectHFS
[ ] Config core.protectNTFS
[ ] Config core.fsmonitor
[ ] Config core.trustMode
[ ] Config core.symlinks
[ ] Config core.bigFileThreshold
[ ] Config core.compression level
[ ] Config core.packedGitWindowSize
[ ] Config core.packedGitLimit
[ ] Config core.deltaBaseCacheLimit
[ ] Config core.preloadIndex
[ ] Config core.untrackedCache
[ ] Config core.quotePath
[ ] Config core.repositoryFormatVersion
[ ] Config core.bare
[ ] Config core.worktree
[ ] Config core.logAllRefUpdates
[ ] Config core.abbrev
[ ] Config core.commentChar
[ ] Config core.fsync
[ ] Config core.fsyncMethod
[ ] Config core.delta
[ ] Config core.multiPackIndex
[ ] Config remote.*.url
[ ] Config remote.*.pushurl
[ ] Config remote.*.fetch refspec
[ ] Config remote.*.push refspec
[ ] Config remote.*.mirror
[ ] Config remote.*.proxy
[ ] Config remote.*.proxyAuthMethod
[ ] Config remote.*.promisor
[ ] Config branch.*.remote
[ ] Config branch.*.merge
[ ] Config branch.*.rebase
[ ] Config branch.*.description
[ ] Config init.defaultBranch
[ ] Config init.templateDir
[ ] Config protocol.version
[ ] Config http.proxy
[ ] Config http.sslVerify
[ ] Config http.sslCert
[ ] Config http.sslKey
[ ] Config http.sslCAInfo
[ ] Config http.postBuffer
[ ] Config http.lowSpeedLimit
[ ] Config http.lowSpeedTime
[ ] Config http.maxRequests
[ ] Config http.cookieFile
[ ] Config http.userAgent
[ ] Config ssh.variant
[ ] Config ssh.command
[ ] Config diff.algorithm
[ ] Config diff.renames
[ ] Config diff.context lines
[ ] Config diff.external tool
[ ] Config merge.tool
[ ] Config merge.conflictStyle
[ ] Config merge.ff
[ ] Config merge.verifySignatures
[ ] Config merge.log
[ ] Config pull.ff
[ ] Config pull.rebase
[ ] Config push.default
[ ] Config push.followTags
[ ] Config push.gpgSign
[ ] Config color.ui
[ ] Config color.status
[ ] Config color.diff
[ ] Config color.branch
[ ] Config credential.helper
[ ] Config credential.username
[ ] Config credential.useHttpPath
[ ] Config filter.*.clean
[ ] Config filter.*.smudge
[ ] Config filter.*.required
[ ] Config filter.*.lfs
[ ] Config alias.*
[ ] Config include.path
[ ] Config includeIf.gitdir
[ ] Config includeIf.branch
[ ] Config safe.directory
[ ] Config core.hooksPath
[ ] Config core.sshCommand
[ ] Config core.alternateObjectDirectories
[ ] Config core.gc config
[ ] Config core.delta config
[ ] Config core.excludesFile
[ ] Config core.attributesFile
[ ] Config core.whitespace
[ ] Config core.fsmonitor config

### Identity — User settings

[ ] Identity command — get/set/delete/list
[ ] Identity update — change user identity
[ ] Identity delete — remove identity
[ ] Identity reset — restore defaults
[ ] Identity clone — copy from existing
[ ] Identity compare — compare two identities
[ ] GetUsername — current user name
[ ] SetUsername — change user name
[ ] IsConfigured — check identity exists
[ ] IsValid — check identity valid
[ ] ValidateUsername — reject empty
[ ] Reject empty usernames — minimum length check
[ ] Trim whitespace — sanitize input
[ ] Enforce maximum username length — limit
[ ] Reject invalid characters — alphanumeric + safe chars
[ ] Validate configuration directory — permissions check
[ ] Initialize configuration directory — create if missing
[ ] Verify configuration — valid identity file
[ ] Reset configuration — clear identity
[ ] Migrate configuration versions — upgrade format
[ ] Backup configuration — before modify
[ ] Restore configuration — from backup
[ ] Load repository identity — per-repo identity
[ ] Save repository identity — per-repo identity file
[ ] Resolve effective identity — repo overrides global
[ ] Check repository identity exists — has repo-level
[ ] Remove repository identity — clear repo-level
[ ] Export identity — portable format
[ ] Import identity — load portable format
[ ] Serialize — binary/text encoding
[ ] Deserialize — parse encoding
[ ] Check identity file exists — file probe
[ ] Get identity file path — path resolution
[ ] Get configuration directory — config dir location
[ ] List configuration files — enumerate configs
[ ] Atomic configuration writes — temp + rename
[ ] Backup before overwrite — safety copy
[ ] Recover interrupted writes — crash recovery
[ ] Restrict configuration file permissions — 0600
[ ] Validate ownership — owner-only access (optional)
[ ] Default identity values — fallback
[ ] Pretty-print identity — human readable
[ ] Equality comparison — two identities equal

## 19. DIFF ENGINE

[ ] Diff type — per-file result
[ ] Diff hunk — contiguous changes
[ ] Diff line — added/removed/context
[ ] Diff line numbers — old/new
[ ] Diff stat — files, insertions, deletions
[ ] Diff shortstat — compact
[ ] Diff numstat — machine format
[ ] Diff raw — raw format
[ ] Diff name-only — file list
[ ] Diff name-status — M/A/D/R/C
[ ] Myers diff O(ND)
[ ] Histogram diff
[ ] Patience diff
[ ] Minimal diff
[ ] Blob vs blob diff
[ ] Tree vs tree diff
[ ] Tree vs working tree
[ ] Commit vs commit
[ ] Diff filter — added/removed/modified
[ ] Diff filter by path
[ ] Diff rename detection — similarity
[ ] Diff rename threshold (default 50%)
[ ] Diff copy detection
[ ] Diff copy threshold
[ ] Diff break detection — rewrite break
[ ] Diff rename limit — max comparisons
[ ] Diff binary detection
[ ] Diff binary marker
[ ] Diff binary patch — base85
[ ] Diff external driver
[ ] Diff algorithm override per file
[ ] Diff textconv — convert binary
[ ] Diff textconv cache
[ ] Diff word diff mode
[ ] Diff word diff regex
[ ] Diff color words
[ ] Diff interhunk context
[ ] Diff function context
[ ] Diff context lines (-U)
[ ] Diff compaction heuristic
[ ] Diff indent heuristic
[ ] Diff anchor function names
[ ] Diff ignore whitespace
[ ] Diff ignore blank lines
[ ] Diff ignore cr at eol
[ ] Diff ignore space at eol
[ ] Diff ignore space change
[ ] Diff ignore all space
[ ] Diff submodule changes
[ ] Diff cached — index vs HEAD
[ ] Diff no-index — outside repo
[ ] Diff relative paths
[ ] Diff unified format output
[ ] Diff color output
[ ] Diff raw output
[ ] Diff patch format headers
[ ] Diff hunk header @@ -a,b +c,d @@
[ ] Diff extended headers — rename/copy
[ ] Diff similarity index
[ ] Diff mode change display
[ ] Diff combined — merge conflict
[ ] Diff combined format
[ ] Diff src-prefix (a/)
[ ] Diff dst-prefix (b/)
[ ] Diff no-prefix
[ ] Diff line-prefix
[ ] Diff function regex per language
[ ] Diff driver selection
[ ] Diff driver built-in
[ ] Diff driver custom config
[ ] Diff memory limit for large files
[ ] Diff algorithm fallback on timeout
[ ] Diff inter-hunk gap merging

## 20. MERGE ENGINE

[ ] Merge type — result of merging
[ ] Merge base — common ancestor
[ ] Merge base find — LCA algorithm
[ ] Merge base octopus — multiple ancestors
[ ] Merge base recursive
[ ] Merge base virtual — synthetic
[ ] Three-way merge — ours/theirs/base
[ ] Three-way file merge
[ ] Three-way hunk merge
[ ] Merge trivial — identical files
[ ] Merge auto-resolve non-conflicting
[ ] Merge conflict detection
[ ] Merge conflict markers — <<< === >>>
[ ] Merge conflict style — merge/diff3/zdiff3
[ ] Merge conflict diff3 — show base
[ ] Merge conflict zdiff3
[ ] Merge conflict labels
[ ] Merge conflict size truncation
[ ] Merge rename/rename conflict
[ ] Merge rename/delete conflict
[ ] Merge modify/delete conflict
[ ] Merge file/directory conflict
[ ] Merge symlink conflict
[ ] Merge binary conflict
[ ] Merge submodule conflict
[ ] Merge tree conflict — structure
[ ] Merge tree recursive — subtrees
[ ] Merge strategy resolve
[ ] Merge strategy recursive
[ ] Merge strategy octopus
[ ] Merge strategy ours
[ ] Merge strategy subtree
[ ] Merge strategy custom plugin
[ ] Merge driver per file
[ ] Merge driver built-in
[ ] Merge driver custom config
[ ] Merge driver config merge.*.name
[ ] Merge driver config merge.*.driver
[ ] Merge verbosity levels
[ ] Merge diff algorithm selection
[ ] Merge rename threshold
[ ] Merge rename limit
[ ] Merge directory rename detection
[ ] Merge directory rename conflict
[ ] Merge directory rename info markers
[ ] Merge conflict notification count
[ ] Merge progress bar
[ ] Merge dry-run
[ ] Merge commit creation
[ ] Merge commit auto-message
[ ] Merge commit template
[ ] Merge commit with log
[ ] Merge no-commit
[ ] Merge no-ff — force merge commit
[ ] Merge ff-only
[ ] Merge ff when possible
[ ] Merge squash
[ ] Merge squash commit
[ ] Merge verify signatures
[ ] Merge abort
[ ] Merge continue
[ ] Merge skip
[ ] Merge tool — external tool
[ ] Merge tool config cmd
[ ] Merge tool trust exit code
[ ] Merge tool keep temporaries
[ ] Merge tool keep backup
[ ] Merge tool prompt
[ ] Merge tool diff preview
[ ] Merge tool GUI mode
[ ] Merge tool list available
[ ] Merge in-progress state files
[ ] Merge in-progress detection
[ ] Merge ref update
[ ] Merge reflog entry
[ ] Merge pre-hook
[ ] Merge post-hook
[ ] Merge hook abort on failure
[ ] Merge commit author
[ ] Merge conflict whitespace ignore
[ ] Merge conflict checkout stages 1/2/3
[ ] Merge stage display
[ ] Merge stage file extraction
[ ] Merge recursive strategy options
[ ] Merge recursive patience
[ ] Merge recursive no-renames
[ ] Merge recursive find-renames threshold
[ ] Merge recursive subtree path
[ ] Merge recursive ours auto-resolve
[ ] Merge recursive theirs auto-resolve

## 21. STATUS — File state display

[ ] Status command
[ ] Tracked file state — modified/added/deleted
[ ] Untracked files
[ ] Ignored files
[ ] Staged changes — index vs HEAD
[ ] Unmerged files — conflicts
[ ] Rename detection
[ ] Copy detection
[ ] Mode changes
[ ] Type changes
[ ] Submodule changes
[ ] Current branch info
[ ] Ahead/behind remote
[ ] Divergence info
[ ] Initial commit state
[ ] Rebase in progress indicator
[ ] Merge in progress indicator
[ ] Cherry-pick in progress
[ ] Revert in progress
[ ] Bisect in progress
[ ] Short format
[ ] Porcelain v1 format
[ ] Porcelain v2 format
[ ] Verbose — show diff
[ ] Untracked mode — normal/all/no
[ ] Ignored mode — normal/no-matching
[ ] Column output
[ ] Relative paths
[ ] Rename limit
[ ] Rename detection toggle
[ ] NUL termination
[ ] Color output
[ ] Sort by path
[ ] Group by status
[ ] Status character — M/A/D/R/C/U/?/!
[ ] Index vs worktree status
[ ] Branch color
[ ] Submodule summary
[ ] Ahead/behind skip toggle
[ ] Ignore submodules config
[ ] Stash count display
[ ] Progress for large trees
[ ] Exclude pattern live update
[ ] Untracked cache
[ ] Untracked cache invalidation
[ ] Sparse checkout info
[ ] Workspace dirty detection
[ ] Conflict file listing

## 22. LOG — Commit history

[ ] Log command
[ ] Log all refs
[ ] Log current branch
[ ] Log specific ref
[ ] Log commit range
[ ] Log path filter
[ ] Log follow renames
[ ] Log max count -n
[ ] Log skip
[ ] Log since date
[ ] Log until date
[ ] Log author filter
[ ] Log committer filter
[ ] Log grep message
[ ] Log pickaxe search
[ ] Log diff per commit
[ ] Log stat per commit
[ ] Log name-only
[ ] Log name-status
[ ] Log raw format
[ ] Log oneline format
[ ] Log short format
[ ] Log medium format (default)
[ ] Log full format
[ ] Log fuller format
[ ] Log reference format
[ ] Log email format
[ ] Log custom format string
[ ] Log tformat with terminator
[ ] Log placeholder %H — full hash
[ ] Log placeholder %h — short hash
[ ] Log placeholder %T — tree hash
[ ] Log placeholder %an — author name
[ ] Log placeholder %ae — author email
[ ] Log placeholder %ad — author date
[ ] Log placeholder %aD — author RFC2822
[ ] Log placeholder %ar — author relative
[ ] Log placeholder %at — author timestamp
[ ] Log placeholder %ai — author ISO 8601
[ ] Log placeholder %cn — committer name
[ ] Log placeholder %ce — committer email
[ ] Log placeholder %cd — committer date
[ ] Log placeholder %cD — committer RFC2822
[ ] Log placeholder %cr — committer relative
[ ] Log placeholder %ct — committer timestamp
[ ] Log placeholder %s — subject
[ ] Log placeholder %f — sanitized subject
[ ] Log placeholder %b — body
[ ] Log placeholder %B — raw body
[ ] Log placeholder %N — notes
[ ] Log placeholder %G? — signature good/bad
[ ] Log placeholder %GS — signer
[ ] Log placeholder %GK — key ID
[ ] Log placeholder %gD — reflog selector
[ ] Log placeholder %gd — short reflog
[ ] Log placeholder %D — decorations
[ ] Log placeholder %d — short decorations
[ ] Log graph — ASCII commit graph
[ ] Log graph branch colors
[ ] Log left-right markers
[ ] Log date format — relative/iso/rfc/short/raw
[ ] Log date order — chronological
[ ] Log topo order — topological
[ ] Log reverse order
[ ] Log ancestry path
[ ] Log first-parent
[ ] Log merge only
[ ] Log no-merges
[ ] Log not — exclude refs
[ ] Log show-signature
[ ] Log decorate refs
[ ] Log decorate auto
[ ] Log mailmap — apply mailmap
[ ] Log output pager
[ ] Log colorized output
[ ] Log expanded ref names

## 23. REBASE — Transplant commits

[ ] Rebase command
[ ] Rebase onto — different base
[ ] Rebase interactive
[ ] Rebase interactive todo — pick/reword/edit/squash/fixup/exec/drop
[ ] Rebase interactive pick — use commit
[ ] Rebase interactive reword — change message
[ ] Rebase interactive edit — amend commit
[ ] Rebase interactive squash — combine
[ ] Rebase interactive fixup — discard message
[ ] Rebase interactive exec — run command
[ ] Rebase interactive drop — omit
[ ] Rebase interactive break — stop
[ ] Rebase interactive label — label position
[ ] Rebase interactive reset — reset to label
[ ] Rebase interactive merge — create merge
[ ] Rebase continue after conflict
[ ] Rebase skip current
[ ] Rebase abort
[ ] Rebase quit
[ ] Rebase apply backend
[ ] Rebase merge backend
[ ] Rebase merge strategy
[ ] Rebase onto base
[ ] Rebase fork-point auto
[ ] Rebase root — entire chain
[ ] Rebase empty — preserve/drop
[ ] Rebase no-ff — force new commits
[ ] Rebase force-rebase
[ ] Rebase committer date
[ ] Rebase ignore date
[ ] Rebase autosquash
[ ] Rebase autostash
[ ] Rebase update-refs
[ ] Rebase ORIG_HEAD
[ ] Rebase state apply dir
[ ] Rebase state merge dir
[ ] Rebase head-name
[ ] Rebase onto name
[ ] Rebase commit count
[ ] Rebase current index
[ ] Rebase interactive done file
[ ] Rebase interactive todo file
[ ] Rebase sign rebased commits
[ ] Rebase conflict detection
[ ] Rebase conflict resolution
[ ] Rebase conflict count
[ ] Rebase diff stat per commit
[ ] Rebase preserve-merges
[ ] Rebase recreate-merges
[ ] Rebase strategy option pass-through
[ ] Rebase ignore-whitespace
[ ] Rebase submodule handling
[ ] Rebase progress bar
[ ] Rebase verbose
[ ] Rebase quiet
[ ] Rebase pre-hook
[ ] Rebase post-hook
[ ] Rebase in-progress detection

## 24. CHERRY-PICK

[ ] Cherry-pick command
[ ] Cherry-pick specific commit
[ ] Cherry-pick range
[ ] Cherry-pick continue
[ ] Cherry-pick skip
[ ] Cherry-pick abort
[ ] Cherry-pick no-commit
[ ] Cherry-pick signoff
[ ] Cherry-pick edit message
[ ] Cherry-pick reflog entry
[ ] Cherry-pick conflict detection
[ ] Cherry-pick conflict resolution
[ ] Cherry-pick parent number
[ ] Cherry-pick mainline
[ ] Cherry-pick strategy
[ ] Cherry-pick strategy option
[ ] Cherry-pick diff algorithm
[ ] Cherry-pick rename detection
[ ] Cherry-pick submodule
[ ] Cherry-pick allow empty
[ ] Cherry-pick empty commit keep/drop
[ ] Cherry-pick state file
[ ] Cherry-pick sequencer batch
[ ] Cherry-pick sign
[ ] Cherry-pick create commit
[ ] Cherry-pick pre-hook
[ ] Cherry-pick post-hook
[ ] Cherry-pick progress
[ ] Cherry-pick verbose
[ ] Cherry-pick quiet
[ ] Cherry-pick in-progress detection
[ ] Cherry-pick reflog message format

## 25. REVERT

[ ] Revert command
[ ] Revert specific commit
[ ] Revert range
[ ] Revert continue
[ ] Revert skip
[ ] Revert abort
[ ] Revert no-commit
[ ] Revert edit message
[ ] Revert mainline
[ ] Revert parent number
[ ] Revert strategy
[ ] Revert signoff
[ ] Revert conflict detection
[ ] Revert conflict resolution
[ ] Revert state file
[ ] Revert sequencer batch
[ ] Revert create commit
[ ] Revert pre-hook
[ ] Revert post-hook
[ ] Revert progress
[ ] Revert verbose
[ ] Revert in-progress detection

## 26. RESET

[ ] Reset command
[ ] Reset soft — move HEAD only
[ ] Reset mixed — move HEAD + index (default)
[ ] Reset hard — move HEAD + index + worktree
[ ] Reset merge — reset merge state
[ ] Reset keep — reset keep local changes
[ ] Reset to commit
[ ] Reset path — specific paths
[ ] Reset path from index — unstage
[ ] Reset path from worktree — discard
[ ] Reset path from HEAD — restore
[ ] Reset path from commit
[ ] Reset ORIG_HEAD set
[ ] Reset reflog entry
[ ] Reset quiet
[ ] Reset verbose
[ ] Reset dry-run
[ ] Reset refresh index
[ ] Reset sparse checkout
[ ] Reset submodule handling
[ ] Reset pre-hook
[ ] Reset post-hook
[ ] Reset in-progress detection

## 27. RESTORE

[ ] Restore command
[ ] Restore worktree — from HEAD
[ ] Restore staged — from HEAD
[ ] Restore from specific commit
[ ] Restore paths — specific files
[ ] Restore pattern — glob
[ ] Restore all tracked
[ ] Restore both staged and worktree
[ ] Restore interactive — per hunk
[ ] Restore patch mode
[ ] Restore progress
[ ] Restore verbose
[ ] Restore quiet
[ ] Restore dry-run
[ ] Restore no-overlay — replace
[ ] Restore overlay — overwrite
[ ] Restore sparse handling
[ ] Restore submodule handling
[ ] Restore staged only
[ ] Restore worktree only
[ ] Restore source ref

## 28. RM

[ ] Rm command
[ ] Rm specific file
[ ] Rm pattern glob
[ ] Rm cached — index only
[ ] Rm force override
[ ] Rm recursive directories
[ ] Rm ignore-unmatch
[ ] Rm dry-run
[ ] Rm quiet
[ ] Rm verbose
[ ] Rm submodule handling
[ ] Rm update index
[ ] Rm backup before removal
[ ] Rm sparse handling
[ ] Rm from HEAD

## 29. MV

[ ] Mv command
[ ] Mv single file
[ ] Mv directory
[ ] Mv force override
[ ] Mv dry-run
[ ] Mv quiet
[ ] Mv verbose
[ ] Mv index update
[ ] Mv worktree move
[ ] Mv rename tracking
[ ] Mv submodule handling
[ ] Mv into existing directory
[ ] Mv overwrite existing
[ ] Mv multiple files — batch
[ ] Mv rename limit

## 30. BRANCH

[ ] Branch command — list/create/delete/rename
[ ] Branch list local
[ ] Branch list all — include remotes
[ ] Branch list merged
[ ] Branch list no-merged
[ ] Branch list verbose — upstream, ahead/behind
[ ] Branch create at HEAD
[ ] Branch create from ref
[ ] Branch create from hash
[ ] Branch create from tag
[ ] Branch create from remote tracking
[ ] Branch delete
[ ] Branch delete force
[ ] Branch delete remote
[ ] Branch rename
[ ] Branch rename force
[ ] Branch copy
[ ] Branch move — alias rename
[ ] Branch set upstream
[ ] Branch unset upstream
[ ] Branch edit description
[ ] Branch show — commit details
[ ] Branch contains check
[ ] Branch merged check
[ ] Branch no-contains
[ ] Branch column output
[ ] Branch sort by key
[ ] Branch sort reverse
[ ] Branch points-at
[ ] Branch format custom
[ ] Branch color output
[ ] Branch abbrev toggle
[ ] Branch recurse submodules
[ ] Branch ignore case
[ ] Branch list pattern
[ ] Branch tracking info
[ ] Branch upstream display
[ ] Branch upstream remote display
[ ] Branch rebase config per-branch
[ ] Branch protection deny delete
[ ] Branch protection deny force push
[ ] Branch config format
[ ] Branch create hook
[ ] Branch delete hook
[ ] Branch rename hook
[ ] Branch prune stale tracking
[ ] Branch prune remote
[ ] Branch detached HEAD warning
[ ] Branch default branch config

## 31. TAG

[ ] Tag command — list/create/delete/verify
[ ] Tag list
[ ] Tag list pattern match
[ ] Tag list sort
[ ] Tag list column output
[ ] Tag list contains commit
[ ] Tag list points-at
[ ] Tag create lightweight
[ ] Tag create annotated
[ ] Tag create from ref
[ ] Tag create from hash
[ ] Tag create replace
[ ] Tag delete
[ ] Tag delete remote
[ ] Tag verify signature
[ ] Tag sign
[ ] Tag message
[ ] Tag message file
[ ] Tag message editor
[ ] Tag force overwrite
[ ] Tag rename
[ ] Tag sort version sort
[ ] Tag column display
[ ] Tag format custom
[ ] Tag color output
[ ] Tag ref — refs/tags/
[ ] Tag object type
[ ] Tag signature inline
[ ] Tag push default
[ ] Tag push follow
[ ] Tag remote operations

## 32. STASH

[ ] Stash command — save/apply/pop/drop/list
[ ] Stash save
[ ] Stash save message
[ ] Stash save untracked
[ ] Stash save all — including ignored
[ ] Stash save keep index
[ ] Stash save patch
[ ] Stash list
[ ] Stash show diff
[ ] Stash show stat
[ ] Stash pop
[ ] Stash pop index restore
[ ] Stash apply
[ ] Stash apply index
[ ] Stash drop
[ ] Stash drop all
[ ] Stash clear
[ ] Stash branch from stash
[ ] Stash ref — refs/stash
[ ] Stash reflog — separate
[ ] Stash count display in status
[ ] Stash conflict detection
[ ] Stash conflict resolution
[ ] Stash index — stash@{N}
[ ] Stash parent — worktree + index
[ ] Stash commit — temporary commit
[ ] Stash create low-level
[ ] Stash store low-level
[ ] Stash push with pathspec
[ ] Stash prune old entries
[ ] Stash drop multiple

## 33. BISECT

[ ] Bisect command — start/bad/good/skip/reset
[ ] Bisect start
[ ] Bisect start bad commit
[ ] Bisect start good commits
[ ] Bisect start terms custom
[ ] Bisect bad
[ ] Bisect good
[ ] Bisect skip
[ ] Bisect skip range
[ ] Bisect reset
[ ] Bisect replay from log
[ ] Bisect log display
[ ] Bisect run script
[ ] Bisect visualize
[ ] Bisect state directory — .steria/BISECT_*
[ ] Bisect commit selection — midpoint
[ ] Bisect commit skip heuristic
[ ] Bisect first-parent mode
[ ] Bisect step count remaining
[ ] Bisect progress display
[ ] Bisect output results
[ ] Bisect multiple ranges
[ ] Bisect terms customizable
[ ] Bisect old/new alternative terms

## 34. CLEAN

[ ] Clean command
[ ] Clean dry-run
[ ] Clean force
[ ] Clean interactive — confirm each
[ ] Clean directories
[ ] Clean excluded — ignored too
[ ] Clean exclude pattern keep
[ ] Clean quiet
[ ] Clean verbose
[ ] Clean pathspec
[ ] Clean mode — files/dirs
[ ] Clean submodule handling
[ ] Clean backup before remove
[ ] Clean sparse handling

## 35. IGNORE — .steriaignore

[ ] Ignore pattern syntax — glob
[ ] Ignore comment — # prefix
[ ] Ignore negation — ! prefix
[ ] Ignore directory only — trailing /
[ ] Ignore recursive — no / = recursive
[ ] Ignore anchored — leading / = root
[ ] Ignore character class — [abc]
[ ] Ignore wildcard — *
[ ] Ignore single char — ?
[ ] Ignore escape — \ prefix
[ ] Ignore double star — **
[ ] Ignore file — .steriaignore
[ ] Ignore per-directory — search up
[ ] Ignore global — ~/.config/steria/ignore
[ ] Ignore local — .steria/info/exclude
[ ] Ignore precedence — most specific wins
[ ] Ignore negation priority — override
[ ] Ignore pattern compile — to regex
[ ] Ignore pattern match — filepath
[ ] Ignore pattern cache — LRU
[ ] Ignore pattern limit
[ ] Ignore nested repos
[ ] Ignore skip-worktree bit
[ ] Ignore assume-unchanged bit
[ ] Ignore status — show ignored
[ ] Ignore check command — test path
[ ] Ignore rule source display
[ ] Ignore debug mode — explain decision

## 36. ATTRIBUTES

[ ] Attributes file — .steriaattributes
[ ] Attributes per-directory precedence
[ ] Attributes global file
[ ] Attributes format — pattern attr=value
[ ] Attributes text/binary
[ ] Attributes crlf — line ending
[ ] Attributes eol — specified ending
[ ] Attributes diff algorithm
[ ] Attributes merge driver
[ ] Attributes filter — clean/smudge
[ ] Attributes export-ignore
[ ] Attributes export-subst
[ ] Attributes language detection
[ ] Attributes working-tree-encoding
[ ] Attributes whitespace rules
[ ] Whitespace tab-in-indent
[ ] Whitespace tabwidth
[ ] Whitespace blank-at-eol
[ ] Whitespace blank-at-eof
[ ] Whitespace space-before-tab
[ ] Whitespace indent-with-non-tab
[ ] Whitespace trailing-space
[ ] Whitespace cr-at-eol
[ ] Attributes macro definition
[ ] Attributes unset
[ ] Attributes unspecified
[ ] Attributes query command
[ ] Attributes debug explain
[ ] Attributes binary — no diff/merge

## 37. HOOKS

[ ] Hook directory — .steria/hooks/
[ ] Hook execution — exec, check exit code
[ ] Hook timeout — kill after duration
[ ] Hook environment — STERIA_DIR vars
[ ] Hook stdin pipe
[ ] Hook stdout capture
[ ] Hook stderr passthrough
[ ] Hook abort — non-zero exit aborts
[ ] Hook allow abort config
[ ] Pre-done hook
[ ] Post-done hook
[ ] Pre-choose hook
[ ] Post-choose hook
[ ] Pre-merge hook
[ ] Post-merge hook
[ ] Pre-rebase hook
[ ] Post-rebase hook
[ ] Pre-push hook
[ ] Post-push hook
[ ] Pre-pull hook
[ ] Post-pull hook
[ ] Pre-fetch hook
[ ] Post-fetch hook
[ ] Pre-receive hook (server)
[ ] Post-receive hook (server)
[ ] Update hook (server per-ref)
[ ] Post-update hook (server)
[ ] Reference-transaction hook
[ ] Prepare-commit-msg hook
[ ] Commit-msg hook — validate message
[ ] Applypatch-msg hook
[ ] Pre-applypatch hook
[ ] Post-applypatch hook
[ ] Fsmonitor hook
[ ] Script format — executable bit + shebang
[ ] Default sample hooks dir
[ ] Hook disable — core.hooksPath
[ ] Hook disable — no-verify flag
[ ] Hook external path config
[ ] Hook parallel execution
[ ] Hook serial execution
[ ] Hook chaining propagation
[ ] Hook security — no shell injection

## 38. REMOTE — Remote management

[ ] Remote command — add/rename/remove/list
[ ] Remote add URL
[ ] Remote add fetch refspec default
[ ] Remote add push refspec default
[ ] Remote add mirror mode
[ ] Remote add no-tags
[ ] Remote add tag option
[ ] Remote rename
[ ] Remote rename update fetch refs
[ ] Remote rename update push refs
[ ] Remote rename update config
[ ] Remote remove
[ ] Remote remove prune refs
[ ] Remote set-url
[ ] Remote set-url push
[ ] Remote set-url add multiple
[ ] Remote set-url delete
[ ] Remote show info
[ ] Remote show branches
[ ] Remote show HEAD
[ ] Remote prune stale
[ ] Remote update — fetch all
[ ] Remote update prune
[ ] Remote get-url
[ ] Remote list
[ ] Remote list verbose
[ ] Remote origin default
[ ] Remote tracking — refs/remotes/
[ ] Remote HEAD detection
[ ] Remote mirror fetch
[ ] Remote mirror push
[ ] Remote sync
[ ] Remote group — multiple remotes

## 39. FETCH

[ ] Fetch command
[ ] Fetch from remote
[ ] Fetch all remotes
[ ] Fetch specific branch
[ ] Fetch all branches
[ ] Fetch tags
[ ] Fetch no-tags
[ ] Fetch force overwrite
[ ] Fetch prune stale
[ ] Fetch prune remote
[ ] Fetch depth — shallow
[ ] Fetch deepen
[ ] Fetch deepen since
[ ] Fetch unshallow
[ ] Fetch update-shallow accept
[ ] Fetch negotiation protocol
[ ] Fetch refspec explicit
[ ] Fetch refmap
[ ] Fetch append — FETCH_HEAD
[ ] Fetch atomic
[ ] Fetch quiet
[ ] Fetch verbose
[ ] Fetch progress
[ ] Fetch server-option
[ ] Fetch submodules recursive
[ ] Fetch submodules on-demand
[ ] Fetch ipv4/6
[ ] Fetch timeout
[ ] Fetch retry
[ ] FETCH_HEAD format
[ ] Fetch auto-gc after
[ ] Fetch write-commit-graph
[ ] Fetch pack handling
[ ] Fetch ref update

## 40. PULL

[ ] Pull command — fetch + merge/rebase
[ ] Pull from remote
[ ] Pull specific branch
[ ] Pull all remotes
[ ] Pull rebase — use rebase
[ ] Pull rebase preserve merges
[ ] Pull rebase interactive
[ ] Pull rebase autostash
[ ] Pull no-rebase — merge
[ ] Pull ff — fast-forward
[ ] Pull ff-only
[ ] Pull no-ff
[ ] Pull squash
[ ] Pull strategy
[ ] Pull strategy option
[ ] Pull verify signatures
[ ] Pull autostash
[ ] Pull depth shallow
[ ] Pull deepen
[ ] Pull unshallow
[ ] Pull server-option
[ ] Pull tags
[ ] Pull no-tags
[ ] Pull progress
[ ] Pull quiet
[ ] Pull verbose
[ ] Pull dry-run
[ ] Pull commit message
[ ] Pull no-commit
[ ] Pull edit message
[ ] Pull cleanup message
[ ] Pull sign merge
[ ] Pull log entries
[ ] Pull signoff
[ ] Pull submodule update
[ ] Pull recurse submodules
[ ] Pull verify-signatures
[ ] Pull summary display
[ ] Pull config per-branch
[ ] Pull into dirty index handling
[ ] Pull conflict handling

## 41. PUSH

[ ] Push command
[ ] Push to remote
[ ] Push branch
[ ] Push all branches
[ ] Push mirror
[ ] Push tags
[ ] Push follow-tags
[ ] Push delete remote ref
[ ] Push force
[ ] Push force-with-lease
[ ] Push force-if-includes
[ ] Push no-verify
[ ] Push dry-run
[ ] Push quiet
[ ] Push verbose
[ ] Push progress
[ ] Push atomic
[ ] Push server-option
[ ] Push signed
[ ] Push push-option
[ ] Push set-upstream
[ ] Push prune
[ ] Push thin pack
[ ] Push no-thin
[ ] Push force-with-lease refspec
[ ] Push ipv4/6
[ ] Push timeout
[ ] Push retry
[ ] Push default mode — simple/current/upstream/matching
[ ] Push output per-ref status
[ ] Push output rejected reasons
[ ] Push refspec explicit mapping
[ ] Push refspec delete — :ref
[ ] Push negotiation
[ ] Push CAS — compare and swap
[ ] Push byte count
[ ] Push speed reporting
[ ] Push pre-hook
[ ] Push transport selection

## 42. TRANSPORT — Protocol layer

[ ] Transport interface — Connect/Disconnect/Push/Pull/Fetch
[ ] Transport smart — capability negotiation
[ ] Transport dumb — static file access
[ ] Transport protocol v1
[ ] Transport protocol v2
[ ] Transport capability list
[ ] Transport capability agent
[ ] Transport capability push-options
[ ] Transport capability fetch
[ ] Transport capability server-option
[ ] Transport capability thin-pack
[ ] Transport capability multi-ack
[ ] Transport capability multi-ack-detailed
[ ] Transport capability side-band
[ ] Transport capability side-band-64k
[ ] Transport capability report-status
[ ] Transport capability report-status-v2
[ ] Transport capability ofs-delta
[ ] Transport capability ref-in-want
[ ] Transport capability promisor-remote
[ ] Transport capability packfile-uris
[ ] Transport capability wait-for-done
[ ] Transport capability shallow
[ ] Transport capability deepen-since
[ ] Transport capability deepen-not
[ ] Transport capability no-progress
[ ] Transport capability include-tag
[ ] Transport ref advertisement
[ ] Transport ref-in-want client request
[ ] Transport fetch request — haves/wants
[ ] Transport fetch response — packfile/ack
[ ] Transport push request — commands + pack
[ ] Transport push response — per-ref status
[ ] Transport side-band multiplexing
[ ] Transport side-band channel — stderr
[ ] Transport side-band channel — pack
[ ] Transport side-band channel — progress
[ ] Transport side-band 64k
[ ] Transport thin pack wire format
[ ] Transport thin pack fix on receive
[ ] Transport packfile URI
[ ] Transport stateless HTTP
[ ] Transport stateful SSH
[ ] Transport gzip encoding
[ ] Transport chunked transfer
[ ] Transport retry — exponential backoff
[ ] Transport redirect following
[ ] Transport redirect limit
[ ] Transport proxy HTTP
[ ] Transport proxy auth
[ ] Transport timeout connect
[ ] Transport timeout data
[ ] Transport timeout total
[ ] Transport user agent header
[ ] Transport cookie support
[ ] Transport cookie file
[ ] Transport extra headers
[ ] Transport low speed abort
[ ] Transport low speed time
[ ] Transport post buffer
[ ] Transport max request size
[ ] Transport SSH command
[ ] Transport SSH variant
[ ] Transport SSH options
[ ] Transport SSH known hosts verify
[ ] Transport HTTP SSL verify
[ ] Transport HTTP SSL cert/key
[ ] Transport HTTP SSL CA
[ ] Transport HTTP auth — basic/digest/bearer
[ ] Transport HTTP auth helper
[ ] Transport pkt-line framing
[ ] Transport pkt-line flush — 0000
[ ] Transport pkt-line delim — 0001
[ ] Transport pkt-line response-end — 0002
[ ] Transport protocol v2 ls-refs
[ ] Transport protocol v2 fetch
[ ] Transport protocol v2 push
[ ] Transport protocol v2 object-info
[ ] Transport protocol v2 bundle-uri
[ ] Transport protocol version detection

## 43. TWO-PHASE SYNC — Steria protocol

[ ] Phase 1: Index sync — POST index
[ ] Phase 1: Server compare index vs store
[ ] Phase 1: Server return missing list
[ ] Phase 2: Client upload missing objects
[ ] Phase 2: PUT /object/:hash single
[ ] Phase 2: POST /objects batch
[ ] Sync init — POST /sync with JSON index
[ ] Sync init request — hash list
[ ] Sync init response — missing hashes
[ ] Sync complete — POST /sync/complete
[ ] Sync push — index + objects
[ ] Sync pull — remote index + delta
[ ] Sync status — ahead/behind
[ ] Sync dry-run — preview
[ ] Sync direction — push/pull/both
[ ] Sync conflict — diverged indices
[ ] Sync conflict resolution — manual/auto
[ ] Sync abort
[ ] Sync resume interrupted
[ ] Sync progress — transfer reporting
[ ] Sync byte count
[ ] Sync object count
[ ] Sync speed rate
[ ] Sync index compression
[ ] Sync index diff — send delta
[ ] Sync index full — entire index
[ ] Sync batch — multiple objects
[ ] Sync batch size config
[ ] Sync batch atomic — all-or-nothing
[ ] Sync verify after sync
[ ] Sync lock
[ ] Sync retry
[ ] Sync timeout
[ ] Sync auth tokens
[ ] Sync sign payloads
[ ] Sync protocol version
[ ] Sync negotiation

## 44. SERVER IMPLEMENTATION

[ ] net/http or fast HTTP server
[ ] Listen address configurable
[ ] Listen port configurable
[ ] TLS cert file
[ ] TLS key file
[ ] TLS auto Let's Encrypt
[ ] TLS mutual auth
[ ] Max connections limit
[ ] Read timeout
[ ] Write timeout
[ ] Idle timeout
[ ] Max header size
[ ] Max body size
[ ] Graceful shutdown
[ ] Route GET /info/refs
[ ] Route POST /sync — sync init
[ ] Route POST /sync/push — push objects
[ ] Route POST /sync/pull — pull objects
[ ] Route GET /object/:hash
[ ] Route PUT /object/:hash
[ ] Route POST /objects batch
[ ] Route GET /repo/info
[ ] Route POST /repo create
[ ] Route DELETE /repo delete
[ ] Route GET /health
[ ] Route GET /metrics
[ ] Router path multiplexer
[ ] Router method dispatch
[ ] Router path parameters
[ ] Router middleware chain
[ ] Middleware logging
[ ] Middleware auth
[ ] Middleware rate limit
[ ] Middleware CORS
[ ] Middleware compression
[ ] Middleware panic recover
[ ] Middleware request ID
[ ] Response JSON encoder
[ ] Response error JSON body
[ ] Response proper HTTP codes
[ ] Response security headers
[ ] Request body read limit
[ ] Request query parser
[ ] Access log format
[ ] Error log format
[ ] Log to file
[ ] Log to stdout
[ ] Rate limit per-IP token bucket
[ ] Rate limit per-token
[ ] Rate limit 429 response
[ ] Rate limit X-RateLimit-* headers
[ ] Auth bearer token verify
[ ] Auth basic auth
[ ] Auth token storage
[ ] Auth token create
[ ] Auth token revoke
[ ] Auth scope read/write/admin
[ ] Auth repo-level permissions
[ ] CORS origin config
[ ] CORS max age
[ ] Compression gzip
[ ] Compression min size
[ ] Compression content types
[ ] Daemon PID file
[ ] Daemon user switch
[ ] Daemon group switch
[ ] Daemon signal handling
[ ] Daemon SIGHUP reload
[ ] Daemon SIGTERM shutdown
[ ] Server storage filesystem backend
[ ] Server storage backend interface
[ ] Server storage backend selection
[ ] Server repo list all
[ ] Server repo create
[ ] Server repo delete
[ ] Server repo description
[ ] Server repo visibility public/private

## 45. CREDENTIAL — Auth storage

[ ] Credential helper interface
[ ] Credential helper built-in cache
[ ] Credential helper built-in store
[ ] Credential helper macOS keychain
[ ] Credential helper Windows wincred
[ ] Credential helper secret service
[ ] Credential helper custom external
[ ] Credential config helper path
[ ] Credential config username
[ ] Credential config useHttpPath
[ ] Credential lookup by URL
[ ] Credential store encrypt file
[ ] Credential store file permissions 0600
[ ] Credential cache in-memory
[ ] Credential cache timeout
[ ] Credential cache clear
[ ] Credential prompt terminal
[ ] Credential prompt hide input
[ ] Credential fill — protocol flow
[ ] Credential approve — store success
[ ] Credential reject — remove stored
[ ] Credential URL parse — host/port/path
[ ] Credential netrc — .netrc support
[ ] Credential netrc machine/login/password

## 46. LOGGING

[ ] Logger interface — Debug/Info/Warn/Error
[ ] Log level debug/info/warn/error/fatal
[ ] Log level configurable
[ ] Log format plain text
[ ] Log format JSON
[ ] Log format key=value
[ ] Log format timestamp
[ ] Log format caller info
[ ] Log format color
[ ] Log output stdout
[ ] Log output stderr
[ ] Log output file
[ ] Log output rotate size
[ ] Log output rotate count
[ ] Log output rotate compress
[ ] Log fields request ID
[ ] Log fields user
[ ] Log fields repo
[ ] Log fields operation
[ ] Log fields duration
[ ] Log sampler rate-limited
[ ] Log context with fields
[ ] Log structured logging
[ ] Log buffer async writes
[ ] Log buffer flush on crash

## 47. ERROR HANDLING

[ ] SteriaError type with code
[ ] Error wrapped chain
[ ] Error code per type
[ ] Error message user-friendly
[ ] Error detail technical
[ ] Error suggestion fix
[ ] Error location file:line
[ ] Error stack trace debug
[ ] Error classification transient/permanent
[ ] Error classification user/system
[ ] Error classification recoverable/fatal
[ ] Error wrapping with %w
[ ] Error format verbosity levels
[ ] Error format color terminal
[ ] Error format JSON machine
[ ] Error sentinel ErrNotFound
[ ] Error sentinel ErrConflict
[ ] Error sentinel ErrNotSteriaProject
[ ] Error sentinel ErrInvalidConfig
[ ] Error sentinel ErrAuthentication
[ ] Error sentinel ErrAuthorization
[ ] Error sentinel ErrObjectCorrupt
[ ] Error sentinel ErrLockFailed
[ ] Error sentinel ErrNotFastForward
[ ] Error sentinel ErrHookFailed
[ ] Error sentinel ErrTimeout
[ ] Defer+recover for panic safety
[ ] Cleanup on error
[ ] Rollback on error
[ ] Auto-retry transient errors
[ ] Exit code mapping per error
[ ] Non-zero exit on failure
[ ] Usage display on bad args
[ ] Error testing injection

## 48. LOCKING

[ ] Lock file .steria/lock
[ ] Lock acquire — O_CREAT|O_EXCL atomic
[ ] Lock release — close + unlink
[ ] Lock timeout — fail after duration
[ ] Lock retry with backoff
[ ] Lock stale detection — PID mismatch
[ ] Lock stale cleanup after timeout
[ ] Lock content — PID/hostname/operation
[ ] Lock non-blocking try
[ ] Lock signal handling — SIGINT release
[ ] Lock defer release
[ ] Lock index separate
[ ] Lock ref separate
[ ] Lock config separate
[ ] Lock shard per directory
[ ] Lock server-side concurrent
[ ] Lock detection — ls *.steria/lock*
[ ] Lock notification error message

## 49. GC — Garbage collection

[ ] GC command — reclaim unused
[ ] GC auto run configurable
[ ] GC auto threshold — loose object count
[ ] GC prune unreachable objects
[ ] GC prune expire — configurable grace
[ ] GC prune now immediate
[ ] GC mark — walk reachable from refs
[ ] GC mark index — active versions
[ ] GC mark head — HEAD commit tree
[ ] GC mark reflog — reflog trees
[ ] GC sweep — delete unreachable loose
[ ] GC sweep verify integrity
[ ] GC pack loose objects
[ ] GC pack window
[ ] GC pack depth
[ ] GC pack bitmap — reachability
[ ] GC pack write
[ ] GC pack delete loose after
[ ] GC repack consolidate packs
[ ] GC repack delta re-delta
[ ] GC repack window merge
[ ] GC reflog expire
[ ] GC reflog expire time
[ ] GC progress bar
[ ] GC quiet
[ ] GC verbose
[ ] GC aggressive — larger window/deeper delta
[ ] GC prune packed after repack
[ ] GC commit graph write
[ ] GC changed path filter

## 50. FSCK — Repository integrity

[ ] Fsck command — check integrity
[ ] Fsck objects — verify all
[ ] Fsck connectivity — DAG check
[ ] Fsck refs — check pointers
[ ] Fsck index — verify
[ ] Fsck reflog — verify
[ ] Fsck config — validate
[ ] Fsck tags annotated format
[ ] Fsck commits message format
[ ] Fsck authors format
[ ] Fsck dates valid
[ ] Fsck tree canonical sort
[ ] Fsck tree duplicates
[ ] Fsck tree valid modes
[ ] Fsck name valid paths
[ ] Fsck dot files check
[ ] Fsck dotgit — .steria in tracked
[ ] Fsck corrupt object — hash mismatch
[ ] Fsck loose corrupt
[ ] Fsck pack corrupt
[ ] Fsck pack CRC mismatch
[ ] Fsck index corrupt
[ ] Fsck recover attempt
[ ] Fsck lost-found unreachable
[ ] Fsck lost-found dir
[ ] Fsck dangling blobs
[ ] Fsck dangling commits
[ ] Fsck dangling trees
[ ] Fsck verbose
[ ] Fsck strict mode
[ ] Fsck error count
[ ] Fsck summary results
[ ] Fsck specific check
[ ] Fsck skip specific check
[ ] Fsck progress
[ ] Fsck missing blob
[ ] Fsck missing tree
[ ] Fsck bad ref target

## 51. SUBMODULE

[ ] Submodule type — external repo ref
[ ] Submodule path on filesystem
[ ] Submodule URL — remote
[ ] Submodule branch tracking
[ ] Submodule add
[ ] Submodule init
[ ] Submodule update — clone/fetch
[ ] Submodule status
[ ] Submodule summary
[ ] Submodule sync — update URL
[ ] Submodule deinit
[ ] Submodule foreach run command
[ ] Submodule set-url
[ ] Submodule set-branch
[ ] Submodule config per-submodule
[ ] Submodule .steriamodules file format
[ ] Submodule nested
[ ] Submodule depth shallow
[ ] Submodule parallel fetch
[ ] Submodule jobs config
[ ] Submodule ignore dirty
[ ] Submodule ignore untracked
[ ] Submodule commit pointer
[ ] Submodule tree entry — gitlink mode
[ ] Submodule diff display
[ ] Submodule merge three-way
[ ] Submodule conflict resolution
[ ] Submodule push order
[ ] Submodule absorb working tree

## 52. WORKTREE

[ ] Worktree command — add/list/prune/lock/unlock/repair
[ ] Worktree add path
[ ] Worktree add branch
[ ] Worktree add commit detached
[ ] Worktree add lock
[ ] Worktree add force
[ ] Worktree add orphan branch
[ ] Worktree add detach
[ ] Worktree list
[ ] Worktree list verbose
[ ] Worktree lock
[ ] Worktree lock reason
[ ] Worktree unlock
[ ] Worktree move
[ ] Worktree prune stale
[ ] Worktree prune dry-run
[ ] Worktree prune expire
[ ] Worktree remove
[ ] Worktree remove force
[ ] Worktree admin dir .steria/worktrees/
[ ] Worktree gitdir file
[ ] Worktree commondir shared
[ ] Worktree ref isolation
[ ] Worktree config isolation
[ ] Worktree index isolation
[ ] Worktree HEAD isolation
[ ] Worktree lock admin dir
[ ] Worktree concurrent access
[ ] Worktree object store sharing

## 53. REFLOG

[ ] Reflog command — show/expire/delete
[ ] Reflog show — all refs
[ ] Reflog show single ref
[ ] Reflog show format
[ ] Reflog expire — prune old
[ ] Reflog expire time config
[ ] Reflog expire unreachable
[ ] Reflog expire dry-run
[ ] Reflog expire verbose
[ ] Reflog delete entry
[ ] Reflog entry — old/new hash, committer, message
[ ] Reflog standard messages per operation
[ ] Reflog HEAD — all HEAD changes
[ ] Reflog branch — per-branch
[ ] Reflog index syntax @{N}
[ ] Reflog date syntax @{date}
[ ] Reflog upstream syntax @{upstream}
[ ] Reflog resolve @{N} to hash
[ ] Reflog file parser
[ ] Reflog file append
[ ] Reflog file per ref
[ ] Reflog file lock
[ ] Reflog limit max entries
[ ] Reflog gc integration

## 54. BLAME

[ ] Blame command — annotate lines
[ ] Blame target file
[ ] Blame from revision
[ ] Blame range lines
[ ] Blame show-name
[ ] Blame show-number
[ ] Blame show-email
[ ] Blame show-age
[ ] Blame line-porcelain
[ ] Blame incremental streaming
[ ] Blame copy detection
[ ] Blame move detection
[ ] Blame ignore-rev
[ ] Blame ignore-revs-file
[ ] Blame ignore-whitespace
[ ] Blame textconv
[ ] Blame algorithm region matching
[ ] Blame pass through
[ ] Blame boundary detection
[ ] Blame first-parent only
[ ] Blame working tree
[ ] Blame binary skip
[ ] Blame large file skip
[ ] Blame progress
[ ] Blame cache diff caching

## 55. GREP

[ ] Grep command — search tracked files
[ ] Grep regex pattern
[ ] Grep fixed string
[ ] Grep extended regex
[ ] Grep ignore-case
[ ] Grep word-regexp
[ ] Grep line-number
[ ] Grep show-function
[ ] Grep recursive
[ ] Grep files-with-matches
[ ] Grep count per file
[ ] Grep name-only
[ ] Grep null delimiter
[ ] Grep context lines
[ ] Grep max-depth
[ ] Grep max-count per file
[ ] Grep and — multiple patterns
[ ] Grep or — multiple patterns
[ ] Grep not — invert
[ ] Grep threads parallel
[ ] Grep at revision
[ ] Grep working tree
[ ] Grep cached index
[ ] Grep no-index outside repo
[ ] Grep recurse-submodules
[ ] Grep binary skip
[ ] Grep textconv
[ ] Grep output color
[ ] Grep pattern cache

## 56. SPARSE CHECKOUT

[ ] Sparse checkout mode — cone/full
[ ] Sparse cone — directories
[ ] Sparse pattern git-style
[ ] Sparse set patterns
[ ] Sparse add pattern
[ ] Sparse list patterns
[ ] Sparse init
[ ] Sparse disable
[ ] Sparse apply — update worktree
[ ] Sparse check-rules test
[ ] Sparse skip-worktree bit
[ ] Sparse status
[ ] Sparse merge handling
[ ] Sparse rebase handling
[ ] Sparse cone depth
[ ] Sparse nested cones
[ ] Sparse partial clone integration
[ ] Sparse progress

## 57. PARTIAL CLONE

[ ] Partial clone — fetch required only
[ ] Filter blob:none
[ ] Filter tree:0
[ ] Filter blob:limit=N
[ ] Filter sparse:oid=
[ ] Filter combine
[ ] Promisor remote config
[ ] Promisor on-demand fetch
[ ] Missing object handling
[ ] Background fetch
[ ] Batch fetch
[ ] Cache local
[ ] Cache eviction
[ ] Prefetch predict
[ ] Server advertise support
[ ] Server serve partial

## 58. FORMAT-PATCH

[ ] Format-patch — generate patches
[ ] Output numbered files
[ ] Output directory
[ ] Commit range
[ ] Cover letter
[ ] Cover letter template
[ ] Patch suffix .patch
[ ] Thread In-Reply-To
[ ] Signoff
[ ] Signature signing
[ ] Attachments MIME
[ ] Inline patch
[ ] Mbox format
[ ] From header
[ ] Date header
[ ] Subject prefix
[ ] Reroll count v2/v3
[ ] Base commit info
[ ] RFC 2822 compliance
[ ] Envelope sender
[ ] SMTP server config
[ ] SMTP port
[ ] SMTP encryption TLS/STARTTLS
[ ] SMTP auth
[ ] Confirm before send
[ ] To/Cc/Bcc headers
[ ] MIME headers

## 59. BUNDLE

[ ] Bundle command — create/list/verify/unbundle
[ ] Bundle create — make file
[ ] Bundle included refs
[ ] Bundle commit range
[ ] Bundle all refs
[ ] Bundle compression zstd
[ ] Bundle progress
[ ] Bundle output path
[ ] Bundle header — prerequisites/refs
[ ] Bundle prerequisite required
[ ] Bundle ref included
[ ] Bundle list refs
[ ] Bundle list prerequisites
[ ] Bundle verify integrity
[ ] Bundle verify prerequisites satisfied
[ ] Bundle unbundle
[ ] Bundle unbundle force
[ ] Bundle version v2/v3
[ ] Bundle thin pack
[ ] Bundle clone from bundle
[ ] Bundle fetch from bundle

## 60. ARCHIVE

[ ] Archive command — tar/zip
[ ] Format tar
[ ] Format zip
[ ] Format tar.gz
[ ] Format tar.zst
[ ] Path prefix in archive
[ ] Output file path
[ ] From commit
[ ] From tree object
[ ] Specific paths
[ ] Current worktree
[ ] Compression level
[ ] Export-ignore attribute
[ ] Export-subst attribute
[ ] Include submodules
[ ] File mode in archive
[ ] Archive remote
[ ] Progress

## 61. SHOW

[ ] Show command — display object
[ ] Show commit format
[ ] Show commit diff
[ ] Show commit stat
[ ] Show tree recursive
[ ] Show tree names
[ ] Show blob raw
[ ] Show tag details
[ ] Show ref pointers
[ ] Show pretty format
[ ] Show abbrev toggle
[ ] Show oneline
[ ] Show quiet
[ ] Show patch
[ ] Show stat
[ ] Show name-only
[ ] Show name-status
[ ] Show color
[ ] Show object type
[ ] Show object size
[ ] Show textconv
[ ] Show pager

## 62. CAT-FILE

[ ] Cat-file — object content
[ ] Cat-file type display
[ ] Cat-file size display
[ ] Cat-file pretty print
[ ] Cat-file textconv
[ ] Cat-file batch mode
[ ] Cat-file batch-check
[ ] Cat-file batch output format
[ ] Cat-file stdin queue
[ ] Cat-file buffer output
[ ] Cat-file stream output
[ ] Cat-file object by hash
[ ] Cat-file with path

## 63. REV-PARSE

[ ] Rev-parse — parse revision
[ ] Rev-parse HEAD
[ ] Rev-parse ref name
[ ] Rev-parse short hash
[ ] Rev-parse verify
[ ] Rev-parse quiet
[ ] Rev-parse show .steria path
[ ] Rev-parse show-toplevel
[ ] Rev-parse show-prefix
[ ] Rev-parse show-cdup
[ ] Rev-parse symbolic-ref
[ ] Rev-parse all refs
[ ] Rev-parse flags parsing
[ ] Rev-parse sq-quote
[ ] Rev-parse spec @{N}
[ ] Rev-parse spec @{-N}
[ ] Rev-parse spec @{upstream}
[ ] Rev-parse spec @{push}
[ ] Rev-parse spec @{/pattern}
[ ] Rev-parse spec :path
[ ] Rev-parse spec ^0/^{tree}/^{commit}
[ ] Rev-parse spec ~N
[ ] Rev-parse spec ^N
[ ] Rev-parse spec .. range
[ ] Rev-parse spec ... symmetric
[ ] Rev-parse spec ^@ parents
[ ] Rev-parse spec ^! commit+parents
[ ] Rev-parse path:hash

## 64. DESCRIBE

[ ] Describe — closest tag
[ ] Describe all refs
[ ] Describe tags only
[ ] Describe contains find
[ ] Describe abbrev length
[ ] Describe match pattern
[ ] Describe exclude pattern
[ ] Describe exact-match
[ ] Describe debug candidates
[ ] Describe long always
[ ] Describe always fallback
[ ] Describe dirty suffix
[ ] Describe first-parent
[ ] Describe output tag-N-gHASH
[ ] Describe search depth

## 65. LARGE FILE STORAGE

[ ] LFS pointer file format
[ ] LFS pointer parse
[ ] LFS pointer create
[ ] LFS smudge filter — download
[ ] LFS clean filter — upload
[ ] LFS batch API — transfer
[ ] LFS verify lock
[ ] LFS lock file
[ ] LFS unlock
[ ] LFS locks list
[ ] LFS transfer adapter — basic
[ ] LFS transfer adapter — SSH
[ ] LFS transfer progress
[ ] LFS retry
[ ] LFS config — url
[ ] LFS config — concurrenttransfers
[ ] LFS config — batch

## 66. PLUMBING COMMANDS

[ ] hash-object — hash a file
[ ] cat-file — display object
[ ] write-tree — write tree
[ ] read-tree — read tree into index
[ ] commit-tree — create commit
[ ] ls-tree — list tree
[ ] ls-files — list index
[ ] update-index — update index
[ ] update-ref — update ref
[ ] rev-list — list commits
[ ] count-objects — count
[ ] verify-pack — verify pack
[ ] unpack-objects — unpack pack
[ ] pack-objects — create pack
[ ] index-pack — build index
[ ] mktag — create tag object
[ ] mktree — create tree object
[ ] diff-files — compare worktree+index
[ ] diff-index — compare index+tree
[ ] diff-tree — compare trees
[ ] for-each-ref — iterate refs
[ ] show-ref — display refs
[ ] symbolic-ref — read/write symref
[ ] check-ref-format — validate name
[ ] prune — prune unreachable
[ ] fsck — check integrity
[ ] gc — garbage collect
[ ] rev-parse — parse revision

## 67. PERFORMANCE

[ ] Object cache LRU
[ ] Delta cache
[ ] Pack memory map
[ ] Preload index
[ ] Untracked cache
[ ] Stat cache
[ ] Commit graph
[ ] Multi-pack index
[ ] Bitmap reachability
[ ] Bloom filter path check
[ ] Thread pool for operations
[ ] Parallel object hash
[ ] Parallel pack compression
[ ] Parallel index diff
[ ] Parallel status scan
[ ] Pipeline read-ahead
[ ] Memory budget config
[ ] Large file threshold
[ ] Timeout per operation
[ ] Progress callback
[ ] Cancel context propagation
[ ] Batch filesystem operations

## 68. SECURITY

[ ] Path traversal prevention
[ ] Command injection prevention
[ ] Shell injection prevention — hooks
[ ] Symlink attack prevention
[ ] .steria dir permissions 0700
[ ] Config file permissions 0600
[ ] SSH key permissions 0600
[ ] Safe directory validation
[ ] Protect HFS+ reserved names
[ ] Protect NTFS reserved names
[ ] Object hash verification on read
[ ] Object hash verification on write
[ ] Signed commit verification
[ ] Signed tag verification
[ ] Push certificate verification
[ ] Reflog tamper detection
[ ] Safe hook execution — timeout
[ ] Hook path validation
[ ] Remote URL validation
[ ] Remote protocol injection
[ ] Transport TLS verification
[ ] Transport host key verification
[ ] Credential storage encryption
[ ] Credential memory clearing
[ ] Token rotation support
[ ] Path allowlist for includes
[ ] Config include chain safety
[ ] Server rate limiting
[ ] Server request size limit
[ ] Server connection limit
[ ] Server auth scope enforcement
[ ] Server path traversal prevention
[ ] Repository isolation
[ ] Safe temp file creation
[ ] Lock file race prevention
[ ] Signal-safe cleanup
[ ] Crash-safe index writes
[ ] Crash-safe ref writes
[ ] Crash-safe object writes
[ ] Fsck integrity enforcement
[ ] Repo repair safety
[ ] GC safety — in-use objects

## 69. MIGRATION

[ ] Index format migration v1→v2
[ ] Pack format migration
[ ] Ref format migration
[ ] Config format migration
[ ] Object format migration
[ ] Commit graph migration
[ ] MIDX migration
[ ] Reflog format migration
[ ] Steria directory layout migration
[ ] Backup before migration
[ ] Rollback on migration failure
[ ] Version compatibility check
[ ] Migration dry-run
[ ] Migration progress
[ ] Migration quiet
[ ] Migration force

## 70. PLATFORM SUPPORT

[ ] Linux — full support
[ ] macOS — full support
[ ] Windows — path handling
[ ] Windows — CRLF handling
[ ] Windows — symlink fallback
[ ] Windows — permission mapping
[ ] Case-insensitive FS detection
[ ] Case-insensitive collision handling
[ ] Unicode normalization
[ ] Long path support Windows
[ ] Path separator normalization
[ ] Temp directory per platform
[ ] Home directory per platform
[ ] Editor detection per platform
[ ] Pager detection
[ ] Shell detection
[ ] Signals per platform
[ ] File locking per platform
[ ] mmap per platform
[ ] FS event monitoring per platform
[ ] ARM vs x86 differences
[ ] 32-bit platform considerations
[ ] Filesystem type detection
[ ] Network FS handling
[ ] NFS safe operations
[ ] SMB safe operations
[ ] FUSE filesystem handling
[ ] Memory-mapped file fallback
[ ] Sparse file support
[ ] Copy-on-write filesystem
[ ] Reflink support
[ ] O_DIRECT support

## 71. CLI — Command line interface

[ ] Global flags — --version, --help
[ ] Global flags — -C <path>, --git-dir
[ ] Global flags --work-tree
[ ] Global flags --bare
[ ] Command dispatch by first arg
[ ] Subcommand tree
[ ] Flag parsing — long/short
[ ] Flag parsing — bool/string/int
[ ] Flag parsing — repeatable
[ ] Flag parsing — negative bool (--no-)
[ ] Argument validation
[ ] Required arguments check
[ ] Argument count check
[ ] Argument type validation
[ ] Usage string generation
[ ] Help text per command
[ ] Help text format consistent
[ ] Man page generation
[ ] Completion bash
[ ] Completion zsh
[ ] Completion fish
[ ] Color detection — terminal
[ ] Color detection — NO_COLOR
[ ] Color detection — TERM
[ ] Pager integration — $PAGER
[ ] Pager fallback — less
[ ] Pager auto-detection
[ ] Output width detection — COLUMNS
[ ] Output width detection — stty
[ ] Progress bar — terminal width
[ ] Progress bar — stderr
[ ] Spinner for indeterminate
[ ] Quiet mode — suppress all
[ ] Verbose mode — debug level
[ ] JSON output mode
[ ] Porcelain output mode
[ ] Exit code — 0 success
[ ] Exit code — 1 general error
[ ] Exit code — 128 bad config
[ ] Exit code — 130 Ctrl+C

## 72. ANSI COLOR

[ ] Color enable — auto/always/never
[ ] Color parse — hex/rgb/named
[ ] Color reset
[ ] Color bold/dim/italic/underline
[ ] Color foreground 8/16/256
[ ] Color background
[ ] Color diff — context/added/removed
[ ] Color branch — current/remote
[ ] Color status — header/file status
[ ] Color log — commit/decorate/graph
[ ] Color grep — match/line/file
[ ] Color interactive — prompt
[ ] Color transport — messages
[ ] Color advance — error/warning
[ ] Color config parsing
[ ] Color config per slot
[ ] Color slot format — fg bg attr
[ ] Color slot validation
[ ] Color slot default
[ ] No-COLOR env var
[ ] Color page — pipe detection
[ ] Color stripping in pipes
[ ] Color cache — parsed values

## 73. INTERACTIVE MODE

[ ] Add -i — interactive staging
[ ] Add -p — patch staging
[ ] Reset -p — patch reset
[ ] Stash -p — patch stash
[ ] Rebase -i — interactive rebase
[ ] Clean -i — interactive clean
[ ] Choose interactive — per-file
[ ] Diff interactive — per-hunk
[ ] Commit interactive — per-file
[ ] Prompt — yes/no confirmation
[ ] Prompt — selection list
[ ] Prompt — text input
[ ] Prompt — password input hidden
[ ] Prompt — editor launch
[ ] Prompt — default values
[ ] Prompt — with/without pager
[ ] Prompt signal handling
[ ] Prompt color

## 74. TEST INFRASTRUCTURE

[ ] Unit test framework
[ ] Integration test framework
[ ] Test helper — temp repo
[ ] Test helper — temp config
[ ] Test helper — seed objects
[ ] Test helper — generate commits
[ ] Test helper — create branches
[ ] Test helper — set identity
[ ] Test fixture — sample repos
[ ] Test fixture — pack files
[ ] Test fixture — signed commits
[ ] Mock transport
[ ] Mock server
[ ] Mock credential helper
[ ] Test coverage target
[ ] Benchmark for hot paths
[ ] Fuzz test for parsers
[ ] Race detection in tests
[ ] Test parallelization
[ ] Test cleanup — remove temps
[ ] Test isolation — no side effects
[ ] Test build flag — integration
[ ] Test for each error path
[ ] Test for each sentinel error
[ ] Test for concurrent access
[ ] Test for signal handling
[ ] Test for crash recovery
[ ] Test for corrupted storage
[ ] Test for large files
[ ] Test for binary files
[ ] Test for Unicode paths
[ ] Test for invalid paths
[ ] Test for symlinks
[ ] Test for submodules
[ ] Test for worktrees
[ ] Test for sparse checkout
[ ] Test for partial clone
[ ] Test for protocol downgrade
[ ] Test for auth failures
[ ] Test for rate limiting
[ ] Test for timeouts
[ ] Test for disk full
[ ] Test for permission denied

## 75. DOCUMENTATION

[ ] Man page — steria(1)
[ ] Man page — steria-watch(1)
[ ] Man page — steria-done(1)
[ ] Man page — steria-choose(1)
[ ] Man page — steria-init(1)
[ ] Man page — steria-clone(1)
[ ] Man page — steria-serve(1)
[ ] Man page — steria-config(1)
[ ] Man page — steria-status(1)
[ ] Man page — steria-log(1)
[ ] Man page — steria-diff(1)
[ ] Man page — steria-branch(1)
[ ] Man page — steria-merge(1)
[ ] Man page — steria-rebase(1)
[ ] Man page — steria-reset(1)
[ ] Man page — steria-stash(1)
[ ] Man page — steria-tag(1)
[ ] Man page — steria-bisect(1)
[ ] Man page — steria-gc(1)
[ ] Man page — steria-fsck(1)
[ ] Man page — steria(7) — design doc
[ ] Man page generation from help
[ ] Man page install path
[ ] README — overview
[ ] README — quick start
[ ] README — commands list
[ ] README — config reference
[ ] README — superposition model
[ ] README — remote protocol
[ ] README — server setup
[ ] README — building from source
[ ] README — contributing guide
[ ] README — license
[ ] Migration guide — from other VCS
[ ] Protocol documentation
[ ] Hooks documentation
[ ] Attributes documentation
[ ] API documentation — Go doc
[ ] Architecture documentation
[ ] Release notes
[ ] Changelog
[ ] Security policy
[ ] Code of conduct
[ ] Build badges
[ ] CLI help text — every command
[ ] Flag documentation — every flag
[ ] Example usage — every command
[ ] Common workflows — documentation
[ ] Best practices — documentation
[ ] Error messages — documented
[ ] Exit code reference
[ ] Config key reference

## 76. BUILD & PACKAGING

[ ] Go build — single binary
[ ] Build tags — platform selection
[ ] Build tags — feature flags
[ ] Build version — ldflags
[ ] Build commit — ldflags
[ ] Build date — ldflags
[ ] CGO_ENABLED=0 — static binary
[ ] Build with CGO for platform features
[ ] Cross-compilation targets
[ ] Linux amd64
[ ] Linux arm64
[ ] macOS amd64
[ ] macOS arm64
[ ] Windows amd64
[ ] FreeBSD amd64
[ ] OpenBSD amd64
[ ] Build script — scripts/build.sh
[ ] Release script — scripts/release.sh
[ ] Package — DEB
[ ] Package — RPM
[ ] Package — APK
[ ] Package — Homebrew
[ ] Package — Scoop
[ ] Package — Docker image
[ ] Dockerfile — alpine
[ ] Dockerfile — distroless
[ ] Docker entrypoint
[ ] Install script — get-steria.sh
[ ] Archive — tar.gz
[ ] Archive — zip
[ ] Checksum — SHA-256
[ ] GPG signature — release
[ ] Version numbering — semver
[ ] Pre-release — alpha/beta/rc
[ ] Build reproducibility
[ ] Vendored dependencies
[ ] Go module proxy
[ ] Minimal build dependencies
[ ] CI — GitHub Actions
[ ] CI — build matrix
[ ] CI — test matrix
[ ] CI — lint step
[ ] CI — vet step
[ ] CI — race detection
[ ] CI — coverage upload
[ ] CI — fuzz testing
[ ] CI — artifact upload
[ ] CI — release automation
[ ] CI — Docker build
[ ] CI — package build
[ ] CI — integration tests

## 77. DEPENDENCY MANAGEMENT

[ ] Go modules — go.mod
[ ] Dependency audit — go vet
[ ] Dependency audit — nancy/govulncheck
[ ] Dependency update process
[ ] Minimal dependency policy
[ ] Stdlib preference
[ ] External dependency justification
[ ] License compliance check
[ ] Dependency vendoring
[ ] Dependency hash verification
[ ] Replace directives — forks
[ ] Indirect dependency pruning

## 78. CODE QUALITY

[ ] gofmt compliance
[ ] go vet compliance
[ ] Staticcheck
[ ] revive linter
[ ] Ineffective assignment check
[ ] Shadow check
[ ] Cyclomatic complexity limit
[ ] Function length limit
[ ] File length limit
[ ] Comment density
[ ] Naming conventions — PascalCase
[ ] Exported function documentation
[ ] Error message conventions
[ ] Log message conventions
[ ] Commit message conventions
[ ] Code review checklist
[ ] Security review checklist
[ ] Performance review checklist

## 79. DEBUG & DIAGNOSE

[ ] Debug mode — --debug flag
[ ] Debug env — STERIA_DEBUG
[ ] Debug trace — all operations
[ ] Debug trace — timing
[ ] Debug trace — memory allocs
[ ] Debug trace — lock acquire/release
[ ] Debug trace — object read/write
[ ] Debug trace — transport bytes
[ ] Debug trace — hook execution
[ ] Debug trace — config loading
[ ] Debug trace — index operations
[ ] Profiling — CPU profile
[ ] Profiling — memory profile
[ ] Profiling — trace profile
[ ] Profiling — mutex profile
[ ] Profiling — block profile
[ ] Profiling output — file
[ ] Profiling output — pprof server
[ ] Stats — STERIA_STATS env
[ ] Stats — per-command counters
[ ] Stats — cache hit rates
[ ] Stats — object counts
[ ] Stats — pack stats
[ ] Stats — transport stats
[ ] Dump — debug dump state
[ ] Dump — environment
[ ] Dump — config active
[ ] Dump — repository layout
[ ] Dump — ref list
[ ] Dump — object stats
[ ] Dump — index contents
[ ] Crash dump — on panic
[ ] Crash dump — repository state
[ ] Crash dump — goroutine stack
[ ] Crash dump — recent logs
[ ] Crash dump — file path
[ ] Crash reporter — optional
[ ] Version info — --version
[ ] Version info — verbose version
[ ] Doctor command — system check
[ ] Doctor — config check
[ ] Doctor — binary verification

## 80. SHELL INTEGRATION

[ ] Bash prompt — __steria_ps1()
[ ] Bash prompt — branch display
[ ] Bash prompt — dirty indicator
[ ] Bash prompt — ahead/behind
[ ] Bash prompt — stash count
[ ] Bash prompt — superposition state
[ ] Bash completion — commands
[ ] Bash completion — branches
[ ] Bash completion — remotes
[ ] Bash completion — tags
[ ] Bash completion — files
[ ] Bash completion — refs
[ ] Bash completion — options
[ ] Zsh prompt integration
[ ] Zsh completion
[ ] Fish prompt integration
[ ] Fish completion
[ ] Prompt performance — no fork bombs
[ ] Prompt config — custom format
[ ] Prompt config — colors
[ ] Prompt config — symbols
[ ] Prompt env vars — STERIA_PS1 vars
[ ] Powerline segment
[ ] Starship integration
[ ] Oh-my-zsh plugin
[ ] Bash-it plugin
[ ] Git-aware prompt compatibility

---

**Total: ~2500+ items across 80 sections**
