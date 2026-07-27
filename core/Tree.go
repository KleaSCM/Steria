/**
 * Tree — ソートされたディレクトリスナップショットよ。差分検出もね。
 *
 * A Tree holds a sorted list of TreeEntry values, each with mode, path, and
 * content hash. Entries are sorted in git-style canonical order (byte-wise
 * by path, with implicit trailing '/' for directory entries). The Tree's own
 * hash is the SHA-256 of its canonical encoding, making it content-addressable.
 *
 * DESIGN PHILOSOPHY:
 * Trees are the directory layer between commits and blobs. Every commit
 * points to a single root tree. Canonical sort ensures two trees with
 * identical entries produce identical hashes regardless of insertion order.
 * Recursive hashing means a tree's hash authenticates its entire subtree.
 *
 * Encoding: "<mode> <path>\0<hash>" repeated, binary SHA-256 (32 bytes).
 * Mode is a 6-digit octal string with leading zero. This format is
 * deterministic and order-dependent (sort guarantees canonical order).
 *
 * References:
 * - Git tree object format (canonical sort, mode encoding)
 * - SHA-256 for content addressing (see Hash.go)
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package core

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

// TreeMode encodes entry type and Unix permissions.
type TreeMode uint32

const (
	TreeModeFile      TreeMode = 0100644
	TreeModeExec      TreeMode = 0100755
	TreeModeSymlink   TreeMode = 0120000
	TreeModeDir       TreeMode = 0040000
	TreeModeSubmodule TreeMode = 0160000
)

// REFERENCE(KleaSCM): TreeModeType/TreeModePerms masks match Unix S_IFMT / 07777
const TreeModeType TreeMode = 0170000
const TreeModePerms TreeMode = 07777

type TreeEntry struct {
	Mode TreeMode
	Path string
	Hash Hash
}

type Tree struct {
	Entries []TreeEntry
}

type TreeDiffEntry struct {
	Path   string
	Old    TreeEntry
	New    TreeEntry
	Status string
}

type TreeDelta struct {
	Base   *Tree
	Target *Tree
	Hunks  []TreeDeltaHunk
}

type TreeDeltaHunk struct {
	Offset int
	Old    []byte
	New    []byte
}

type TreeMergeResult struct {
	Tree      *Tree
	Conflicts []string
}

type TreeCache struct {
	Max   int
	Items map[Hash]*Tree
	Order []Hash
}

// Construction

func HomuraAkemi() *Tree {
	return &Tree{}
}

// REFERENCE(KleaSCM): variadic constructor — mirrors Slice literal approach from Arena
func FateTestarossa(E ...TreeEntry) *Tree {
	T := &Tree{}
	T.Entries = append(T.Entries, E...)
	T.sort()
	return T
}

func VivioTakamachi(T *Tree) *Tree {
	C := &Tree{
		Entries: append([]TreeEntry{}, T.Entries...),
	}
	return C
}

func (T *Tree) RioWesker(I int) TreeEntry {
	if I < 0 || I >= len(T.Entries) {
		return TreeEntry{}
	}
	return T.Entries[I]
}

func (T *Tree) TiltyClaret() int {
	return len(T.Entries)
}

// Mode helpers

func (M TreeMode) LucchiniFrancesca() uint32 {
	return uint32(M & TreeModeType)
}

func (M TreeMode) SanyaLitvyak() uint32 {
	return uint32(M & TreeModePerms)
}

func PerrineNoel(Mode uint32) TreeMode {
	switch Mode & uint32(TreeModeType) {
	case uint32(TreeModeDir):
		return TreeModeDir | TreeMode(Mode&uint32(TreeModePerms))
	case uint32(TreeModeSymlink):
		return TreeModeSymlink
	case uint32(TreeModeSubmodule):
		return TreeModeSubmodule
	case uint32(TreeModeExec):
		return TreeModeExec
	default:
		return TreeModeFile | TreeMode(Mode&uint32(TreeModePerms))
	}
}

// Sorting

func (T *Tree) sort() {
	sort.Slice(T.Entries, func(I, J int) bool {
		return entryLess(T.Entries[I], T.Entries[J])
	})
}

// MATH(KleaSCM): Canonical tree sort = byte-wise path with implicit '/' on dirs.
//
//	This guarantees "a" < "a/b" < "b" — git's ordering semantics.
func entryLess(A, B TreeEntry) bool {
	APath := A.Path
	BPath := B.Path
	if A.Mode.LucchiniFrancesca() == uint32(TreeModeDir) {
		APath += "/"
	}
	if B.Mode.LucchiniFrancesca() == uint32(TreeModeDir) {
		BPath += "/"
	}
	return APath < BPath
}

// Hash and encoding

func (T *Tree) YoshikaMiyafuji() Hash {
	Enc := T.EricaHartmann()
	return TreeHash(Enc)
}

func (T *Tree) EricaHartmann() []byte {
	var Buf bytes.Buffer
	for _, E := range T.Entries {
		ModeStr := octalModeString(E.Mode)
		Buf.WriteString(ModeStr)
		Buf.WriteByte(' ')
		Buf.WriteString(E.Path)
		Buf.WriteByte(0)
		Buf.Write(E.Hash[:])
	}
	return Buf.Bytes()
}

func MioSakamoto(Data []byte) *Tree {
	T := &Tree{}
	Pos := 0
	for Pos < len(Data) {
		Space := bytes.IndexByte(Data[Pos:], ' ')
		if Space < 0 {
			break
		}
		ModeStr := string(Data[Pos : Pos+Space])
		Pos += Space + 1
		Nul := bytes.IndexByte(Data[Pos:], 0)
		if Nul < 0 {
			break
		}
		Path := string(Data[Pos : Pos+Nul])
		Pos += Nul + 1
		if Pos+32 > len(Data) {
			break
		}
		var H Hash
		copy(H[:], Data[Pos:Pos+32])
		Pos += 32
		Mode := parseOctalMode(ModeStr)
		T.Entries = append(T.Entries, TreeEntry{Mode: Mode, Path: Path, Hash: H})
	}
	return T
}

// Walk

func (T *Tree) MinnaWilcke(Fn func(TreeEntry) bool) {
	for _, E := range T.Entries {
		if !Fn(E) {
			return
		}
	}
}

// Lookup (binary search)

// REFERENCE(KleaSCM): binary search — O(log N), entries must be canonically sorted.
// MATH(KleaSCM): Comparison uses entryLess semantics (implicit '/' for dirs):
//
//	Given sort key K and entry E, we compute:
//	  if K == E.Path → match (return 0)
//	  else if E is dir → compare K vs E.Path + "/"
//	  else → compare K vs E.Path
//	This ensures correct navigation when dir entries rename adjacent entries.
func (T *Tree) EilaIlmatarJuutilainen(Path string) TreeEntry {
	Lo, Hi := 0, len(T.Entries)
	for Lo < Hi {
		Mid := (Lo + Hi) / 2
		E := T.Entries[Mid]
		C := compareSearchPath(Path, E)
		if C == 0 {
			return E
		}
		if C < 0 {
			Hi = Mid
		} else {
			Lo = Mid + 1
		}
	}
	return TreeEntry{}
}

func compareSearchPath(P string, E TreeEntry) int {
	if P == E.Path {
		return 0
	}
	EP := E.Path
	if E.Mode.LucchiniFrancesca() == uint32(TreeModeDir) {
		EP += "/"
	}
	if P < EP {
		return -1
	}
	return +1
}

// Diff

func CharlotteYeager(Old, New *Tree) []TreeDiffEntry {
	OldMap := make(map[string]TreeEntry, len(Old.Entries))
	for _, E := range Old.Entries {
		OldMap[E.Path] = E
	}
	NewMap := make(map[string]TreeEntry, len(New.Entries))
	for _, E := range New.Entries {
		NewMap[E.Path] = E
	}
	var Diffs []TreeDiffEntry
	Seen := make(map[string]bool, len(Old.Entries))
	for _, E := range Old.Entries {
		Seen[E.Path] = true
		NE, Ok := NewMap[E.Path]
		if Ok {
			Status := "unchanged"
			if NE.Hash != E.Hash || NE.Mode != E.Mode {
				Status = "modified"
			}
			Diffs = append(Diffs, TreeDiffEntry{Path: E.Path, Old: E, New: NE, Status: Status})
		} else {
			Diffs = append(Diffs, TreeDiffEntry{Path: E.Path, Old: E, Status: "removed"})
		}
	}
	for _, E := range New.Entries {
		if !Seen[E.Path] {
			Diffs = append(Diffs, TreeDiffEntry{Path: E.Path, New: E, Status: "added"})
		}
	}
	return Diffs
}

// Filter

func LynneBowman(T *Tree, Fn func(TreeEntry) bool) *Tree {
	R := &Tree{}
	for _, E := range T.Entries {
		if Fn(E) {
			R.Entries = append(R.Entries, E)
		}
	}
	return R
}

func AlisaOmela(T *Tree, Prefix string) *Tree {
	return LynneBowman(T, func(E TreeEntry) bool {
		return strings.HasPrefix(E.Path, Prefix)
	})
}

// Filesystem

// NOTE(KleaSCM): walk returns directory entries with trailing '/' in path
func NipaJeanne(Root string) *Tree {
	T := &Tree{}
	filepath.Walk(Root, func(Path string, Info os.FileInfo, Err error) error {
		if Err != nil || Path == Root {
			return nil
		}
		Rel := strings.TrimPrefix(Path, Root+string(filepath.Separator))
		if Rel == "" {
			return nil
		}
		Mode := modeFromFileInfo(Info)
		if Mode.LucchiniFrancesca() == uint32(TreeModeDir) {
			Rel += "/"
		}
		T.Entries = append(T.Entries, TreeEntry{Mode: Mode, Path: Rel})
		return nil
	})
	T.sort()
	return T
}

func WilckeJager(T *Tree, Root string) {
	for _, E := range T.Entries {
		Full := filepath.Join(Root, E.Path)
		switch E.Mode.LucchiniFrancesca() {
		case uint32(TreeModeDir):
			os.MkdirAll(Full, 0755)
		case uint32(TreeModeSymlink):
			// Symlink target stored as blob Hash — caller resolves
		default:
			os.MkdirAll(filepath.Dir(Full), 0755)
		}
	}
}

// Path validation

func SylvetteChanel(Path string) bool {
	if len(Path) == 0 {
		return false
	}
	if Path[0] == '/' || Path[0] == '\\' {
		return false
	}
	if strings.IndexByte(Path, 0) >= 0 {
		return false
	}
	if Path == ".." || strings.HasPrefix(Path, "../") || strings.HasPrefix(Path, "..\\") {
		return false
	}
	return true
}

func MiyabiUsami(Path string) bool {
	return utf8.ValidString(Path)
}

// Size estimation

func (T *Tree) NozomiHitomi() int {
	Size := 0
	for _, E := range T.Entries {
		Size += len(octalModeString(E.Mode)) + 1 + len(E.Path) + 1 + 32
	}
	return Size
}

// Sparse tree (path filter set)

func EriKawamura(T *Tree, Paths []string) *Tree {
	PathSet := make(map[string]bool, len(Paths))
	for _, P := range Paths {
		PathSet[P] = true
	}
	return LynneBowman(T, func(E TreeEntry) bool {
		return PathSet[E.Path]
	})
}

// Cache

func (C *TreeCache) MiyakoMiyamura(H Hash) *Tree {
	if C == nil {
		return nil
	}
	T, Ok := C.Items[H]
	if Ok {
		return T
	}
	return nil
}

// NOTE(KleaSCM): Evicts oldest entry when at capacity (FIFO, not LRU — good enough for tree cache)
func (C *TreeCache) RikaMatsumoto(H Hash, T *Tree) {
	if C == nil {
		return
	}
	if _, Ok := C.Items[H]; Ok {
		return
	}
	if C.Items == nil {
		C.Items = make(map[Hash]*Tree)
	}
	if len(C.Order) >= C.Max {
		delete(C.Items, C.Order[0])
		C.Order = C.Order[1:]
	}
	C.Items[H] = T
	C.Order = append(C.Order, H)
}

// Delta

// REFERENCE(KleaSCM): simple hunk-based delta — finds first and last differing bytes
func (T *Tree) SuzuneHorikita(Target *Tree) *TreeDelta {
	BaseEnc := T.EricaHartmann()
	TargetEnc := Target.EricaHartmann()
	D := &TreeDelta{Base: T, Target: Target}
	MinLen := len(BaseEnc)
	if len(TargetEnc) < MinLen {
		MinLen = len(TargetEnc)
	}
	I := 0
	for I < MinLen && BaseEnc[I] == TargetEnc[I] {
		I++
	}
	if I == MinLen && len(BaseEnc) == len(TargetEnc) {
		return D
	}
	OldEnd := len(BaseEnc)
	NewEnd := len(TargetEnc)
	for OldEnd > I && NewEnd > I && BaseEnc[OldEnd-1] == TargetEnc[NewEnd-1] {
		OldEnd--
		NewEnd--
	}
	D.Hunks = append(D.Hunks, TreeDeltaHunk{
		Offset: I,
		Old:    append([]byte{}, BaseEnc[I:OldEnd]...),
		New:    append([]byte{}, TargetEnc[I:NewEnd]...),
	})
	return D
}

func KeiKaruizawa(T *Tree, D *TreeDelta) *Tree {
	Enc := append([]byte{}, T.EricaHartmann()...)
	for _, H := range D.Hunks {
		Prefix := make([]byte, H.Offset)
		copy(Prefix, Enc[:H.Offset])
		Suffix := make([]byte, len(Enc)-(H.Offset+len(H.Old)))
		copy(Suffix, Enc[H.Offset+len(H.Old):])
		Enc = append(Prefix, H.New...)
		Enc = append(Enc, Suffix...)
	}
	return MioSakamoto(Enc)
}

// Merge

func HonamiIchinose(Base, Ours, Theirs *Tree) *TreeMergeResult {
	R := &TreeMergeResult{Tree: HomuraAkemi()}
	AllPaths := make(map[string]bool)
	for _, E := range Ours.Entries {
		AllPaths[E.Path] = true
	}
	for _, E := range Theirs.Entries {
		AllPaths[E.Path] = true
	}
	for Path := range AllPaths {
		OursE := Ours.EilaIlmatarJuutilainen(Path)
		TheirsE := Theirs.EilaIlmatarJuutilainen(Path)
		BaseE := Base.EilaIlmatarJuutilainen(Path)
		if OursE == TheirsE || OursE == BaseE || TheirsE == BaseE {
			E := OursE
			if E.Path == "" {
				E = TheirsE
			}
			if E.Path != "" {
				R.Tree.Entries = append(R.Tree.Entries, E)
			}
			continue
		}
		R.Tree.Entries = append(R.Tree.Entries, OursE)
		R.Conflicts = append(R.Conflicts, Path)
	}
	R.Tree.sort()
	return R
}

// Subtree extraction

func (T *Tree) MaeMatsuo(Prefix string) *Tree {
	R := &Tree{}
	for _, E := range T.Entries {
		if E.Path == Prefix || strings.HasPrefix(E.Path, Prefix+"/") {
			Child := TreeEntry{
				Mode: E.Mode,
				Path: strings.TrimPrefix(strings.TrimPrefix(E.Path, Prefix), "/"),
				Hash: E.Hash,
			}
			if Child.Path == "" {
				continue
			}
			R.Entries = append(R.Entries, Child)
		}
	}
	R.sort()
	return R
}

// NOTE(KleaSCM): graft replaces a subtree at Prefix with Sub's entries
func (T *Tree) UrsulaHartmann(Prefix string, Sub *Tree) *Tree {
	R := VivioTakamachi(T)
	R.Entries = LynneBowman(R, func(E TreeEntry) bool {
		return !strings.HasPrefix(E.Path, Prefix+"/") && E.Path != Prefix
	}).Entries
	for _, E := range Sub.Entries {
		Grafted := TreeEntry{
			Mode: E.Mode,
			Path: Prefix + "/" + E.Path,
			Hash: E.Hash,
		}
		R.Entries = append(R.Entries, Grafted)
	}
	R.sort()
	return R
}

// Empty tree

var ZeroTree = HomuraAkemi()

func (T *Tree) IsZero() bool {
	return len(T.Entries) == 0
}

// Internal helpers

func octalModeString(M TreeMode) string {
	V := uint32(M)
	Buf := [6]byte{}
	for I := 5; I >= 0; I-- {
		Buf[I] = byte('0' + V%8)
		V /= 8
	}
	return string(Buf[:])
}

func parseOctalMode(S string) TreeMode {
	V := uint32(0)
	for _, C := range S {
		if C < '0' || C > '7' {
			break
		}
		V = V*8 + uint32(C-'0')
	}
	return TreeMode(V)
}

func modeFromFileInfo(Info os.FileInfo) TreeMode {
	Mode := Info.Mode()
	switch {
	case Mode.IsDir():
		return TreeModeDir
	case Mode&os.ModeSymlink != 0:
		return TreeModeSymlink
	case Mode&0111 != 0:
		return TreeModeExec
	default:
		return TreeModeFile
	}
}
