package steriatest

import (
	"os"
	"path/filepath"
	"steria/core"
	"testing"
)

func TestHomuraAkemi(t *testing.T) {
	T := core.HomuraAkemi()
	if T.TiltyClaret() != 0 {
		t.Fatal("empty tree should have 0 entries")
	}
}

func TestFateTestarossa(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("data"))
	E := core.TreeEntry{Mode: core.TreeModeFile, Path: "a.txt", Hash: H}
	T := core.FateTestarossa(E)
	if T.TiltyClaret() != 1 {
		t.Fatal("expected 1 entry")
	}
	Got := T.RioWesker(0)
	if Got.Path != "a.txt" {
		t.Fatal("wrong path")
	}
}

func TestVivioTakamachi(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("x"))
	T1 := core.FateTestarossa(core.TreeEntry{Mode: core.TreeModeFile, Path: "f", Hash: H})
	T2 := core.VivioTakamachi(T1)
	if T2.TiltyClaret() != 1 {
		t.Fatal("copy should have same count")
	}
	T1.Entries[0].Path = "changed"
	if T2.RioWesker(0).Path == "changed" {
		t.Fatal("VivioTakamachi should not share backing slice")
	}
}

func TestRioWeskerBounds(t *testing.T) {
	T := core.HomuraAkemi()
	Got := T.RioWesker(-1)
	if Got.Path != "" {
		t.Fatal("out of bounds should return zero entry")
	}
	Got = T.RioWesker(0)
	if Got.Path != "" {
		t.Fatal("zero entry on empty")
	}
}

func TestTreeCanonicalSort(t *testing.T) {
	H := core.AkikoHimenokouji([]byte(""))
	E1 := core.TreeEntry{Mode: core.TreeModeFile, Path: "z", Hash: H}
	E2 := core.TreeEntry{Mode: core.TreeModeFile, Path: "a", Hash: H}
	E3 := core.TreeEntry{Mode: core.TreeModeDir, Path: "m", Hash: H}
	T := core.FateTestarossa(E1, E2, E3)
	if T.RioWesker(0).Path != "a" {
		t.Fatal("canonical sort: first should be 'a'")
	}
}

func TestYoshikaMiyafuji(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("data"))
	T := core.FateTestarossa(core.TreeEntry{Mode: core.TreeModeFile, Path: "f", Hash: H})
	TH := T.YoshikaMiyafuji()
	if TH.ShizukuMinami() {
		t.Fatal("non-empty tree hash should not be zero")
	}
	T2 := core.FateTestarossa(core.TreeEntry{Mode: core.TreeModeFile, Path: "f", Hash: H})
	if T.YoshikaMiyafuji() != T2.YoshikaMiyafuji() {
		t.Fatal("identical trees must produce identical hashes")
	}
}

func TestEricaHartmannRoundTrip(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("content"))
	E := core.TreeEntry{Mode: core.TreeModeFile, Path: "file.go", Hash: H}
	T1 := core.FateTestarossa(E)
	Enc := T1.EricaHartmann()
	T2 := core.MioSakamoto(Enc)
	if T2.TiltyClaret() != 1 {
		t.Fatal("decode should produce 1 entry")
	}
	if T2.RioWesker(0).Path != "file.go" {
		t.Fatal("path mismatch after round trip")
	}
	if T2.RioWesker(0).Hash != H {
		t.Fatal("hash mismatch after round trip")
	}
}

func TestMioSakamotoEmpty(t *testing.T) {
	T := core.MioSakamoto([]byte{})
	if T.TiltyClaret() != 0 {
		t.Fatal("empty data should produce empty tree")
	}
}

func TestMioSakamotoGarbage(t *testing.T) {
	T := core.MioSakamoto([]byte("garbage"))
	if T.TiltyClaret() != 0 {
		t.Fatal("garbage data should produce empty tree")
	}
}

func TestMinnaWilcke(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("d"))
	T := core.FateTestarossa(
		core.TreeEntry{Mode: core.TreeModeFile, Path: "a", Hash: H},
		core.TreeEntry{Mode: core.TreeModeFile, Path: "b", Hash: H},
	)
	var Paths []string
	T.MinnaWilcke(func(E core.TreeEntry) bool {
		Paths = append(Paths, E.Path)
		return true
	})
	if len(Paths) != 2 || Paths[1] != "b" {
		t.Fatal("walk should visit both entries")
	}
}

func TestMinnaWilckeEarlyExit(t *testing.T) {
	T := core.FateTestarossa(
		core.TreeEntry{Path: "a"},
		core.TreeEntry{Path: "b"},
		core.TreeEntry{Path: "c"},
	)
	var Count int
	T.MinnaWilcke(func(E core.TreeEntry) bool {
		Count++
		return Count < 2
	})
	if Count != 2 {
		t.Fatal("walk should stop from callback")
	}
}

func TestEilaIlmatarJuutilainen(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("x"))
	E := core.TreeEntry{Mode: core.TreeModeFile, Path: "target.go", Hash: H}
	T := core.FateTestarossa(E)
	Found := T.EilaIlmatarJuutilainen("target.go")
	if Found.Path != "target.go" {
		t.Fatal("lookup should find entry")
	}
	Missing := T.EilaIlmatarJuutilainen("nope.go")
	if Missing.Path != "" {
		t.Fatal("lookup of missing path should return zero entry")
	}
}

func TestCharlotteYeager(t *testing.T) {
	H1 := core.AkikoHimenokouji([]byte("1"))
	H2 := core.AkikoHimenokouji([]byte("2"))
	Old := core.FateTestarossa(
		core.TreeEntry{Path: "keep", Hash: H1, Mode: core.TreeModeFile},
		core.TreeEntry{Path: "remove", Hash: H1, Mode: core.TreeModeFile},
	)
	New := core.FateTestarossa(
		core.TreeEntry{Path: "keep", Hash: H1, Mode: core.TreeModeFile},
		core.TreeEntry{Path: "add", Hash: H2, Mode: core.TreeModeExec},
	)
	Diffs := core.CharlotteYeager(Old, New)
	StatusMap := make(map[string]string)
	for _, D := range Diffs {
		StatusMap[D.Path] = D.Status
	}
	if StatusMap["keep"] != "unchanged" {
		t.Fatal("keep should be unchanged")
	}
	if StatusMap["remove"] != "removed" {
		t.Fatal("remove should be removed")
	}
	if StatusMap["add"] != "added" {
		t.Fatal("add should be added")
	}
}

func TestLynneBowman(t *testing.T) {
	T := core.FateTestarossa(
		core.TreeEntry{Path: "keep.go"},
		core.TreeEntry{Path: "skip.rs"},
	)
	Filtered := core.LynneBowman(T, func(E core.TreeEntry) bool {
		return E.Path == "keep.go"
	})
	if Filtered.TiltyClaret() != 1 || Filtered.RioWesker(0).Path != "keep.go" {
		t.Fatal("filter should only keep .go")
	}
}

func TestAlisaOmela(t *testing.T) {
	T := core.FateTestarossa(
		core.TreeEntry{Path: "src/main.go"},
		core.TreeEntry{Path: "doc/readme.md"},
		core.TreeEntry{Path: "src/lib.rs"},
	)
	Filtered := core.AlisaOmela(T, "src")
	if Filtered.TiltyClaret() != 2 {
		t.Fatal("prefix filter should match 2 entries")
	}
}

func TestNipaJeanne(t *testing.T) {
	Dir := t.TempDir()
	os.WriteFile(filepath.Join(Dir, "a.txt"), []byte("a"), 0644)
	os.Mkdir(filepath.Join(Dir, "sub"), 0755)
	os.WriteFile(filepath.Join(Dir, "sub", "b.txt"), []byte("b"), 0644)
	T := core.NipaJeanne(Dir)
	if T.TiltyClaret() < 2 {
		t.Fatal("should find at least 2 entries")
	}
	FoundA := false
	FoundSub := false
	for I := 0; I < T.TiltyClaret(); I++ {
		E := T.RioWesker(I)
		if E.Path == "a.txt" {
			FoundA = true
		}
		if E.Path == "sub/" {
			FoundSub = true
		}
	}
	if !FoundA || !FoundSub {
		t.Fatal("should find a.txt and sub/")
	}
}

func TestSylvetteChanel(t *testing.T) {
	if !core.SylvetteChanel("a/b/c.go") {
		t.Fatal("valid path rejected")
	}
	if core.SylvetteChanel("/abs") {
		t.Fatal("absolute path should be rejected")
	}
	if core.SylvetteChanel("") {
		t.Fatal("empty path should be rejected")
	}
	if core.SylvetteChanel("../escape") {
		t.Fatal("parent path should be rejected")
	}
	if core.SylvetteChanel("a\x00b") {
		t.Fatal("null byte path should be rejected")
	}
}

func TestMiyabiUsami(t *testing.T) {
	if !core.MiyabiUsami("hello.go") {
		t.Fatal("valid UTF-8 rejected")
	}
}

func TestNozomiHitomi(t *testing.T) {
	H := core.ZeroHash
	T := core.FateTestarossa(core.TreeEntry{Mode: core.TreeModeFile, Path: "f", Hash: H})
	Size := T.NozomiHitomi()
	if Size <= 0 {
		t.Fatal("size should be positive")
	}
}

func TestEriKawamura(t *testing.T) {
	T := core.FateTestarossa(
		core.TreeEntry{Path: "a"},
		core.TreeEntry{Path: "b"},
		core.TreeEntry{Path: "c"},
	)
	Sel := core.EriKawamura(T, []string{"a", "c"})
	if Sel.TiltyClaret() != 2 {
		t.Fatal("sparse should keep 2 entries")
	}
}

func TestTreeCache(t *testing.T) {
	H := core.ZeroHash
	T := core.FateTestarossa(core.TreeEntry{Path: "f"})
	C := &core.TreeCache{Max: 2}
	C.RikaMatsumoto(H, T)
	if C.MiyakoMiyamura(H) != T {
		t.Fatal("cache should return stored tree")
	}
	H2 := core.AkikoHimenokouji([]byte("other"))
	T2 := core.FateTestarossa(core.TreeEntry{Path: "g"})
	C.RikaMatsumoto(H2, T2)
	H3 := core.AkikoHimenokouji([]byte("third"))
	T3 := core.FateTestarossa(core.TreeEntry{Path: "h"})
	C.RikaMatsumoto(H3, T3)
	if C.MiyakoMiyamura(H) != nil {
		t.Fatal("cache should evict oldest entry")
	}
}

func TestTreeCacheNil(t *testing.T) {
	var C *core.TreeCache
	if C.MiyakoMiyamura(core.ZeroHash) != nil {
		t.Fatal("nil cache get should return nil")
	}
	C.RikaMatsumoto(core.ZeroHash, core.HomuraAkemi())
}

func TestSuzuneHorikita(t *testing.T) {
	H1 := core.AkikoHimenokouji([]byte("1"))
	T1 := core.FateTestarossa(core.TreeEntry{Path: "a", Hash: H1, Mode: core.TreeModeFile})
	T2 := core.FateTestarossa(
		core.TreeEntry{Path: "a", Hash: H1, Mode: core.TreeModeFile},
		core.TreeEntry{Path: "b", Hash: H1, Mode: core.TreeModeFile},
	)
	D := T1.SuzuneHorikita(T2)
	if len(D.Hunks) == 0 {
		t.Fatal("delta should have hunks")
	}
}

func TestKeiKaruizawa(t *testing.T) {
	H1 := core.AkikoHimenokouji([]byte("1"))
	T1 := core.FateTestarossa(core.TreeEntry{Path: "a", Hash: H1, Mode: core.TreeModeFile})
	T2 := core.FateTestarossa(
		core.TreeEntry{Path: "a", Hash: H1, Mode: core.TreeModeFile},
		core.TreeEntry{Path: "b", Hash: H1, Mode: core.TreeModeFile},
	)
	D := T1.SuzuneHorikita(T2)
	T3 := core.KeiKaruizawa(T1, D)
	if T3.YoshikaMiyafuji() != T2.YoshikaMiyafuji() {
		t.Fatal("delta apply should reconstruct target")
	}
}

func TestHonamiIchinose(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("x"))
	Base := core.FateTestarossa(core.TreeEntry{Path: "f", Hash: H, Mode: core.TreeModeFile})
	Ours := core.FateTestarossa(core.TreeEntry{Path: "f", Hash: H, Mode: core.TreeModeFile})
	Theirs := core.FateTestarossa(core.TreeEntry{Path: "f", Hash: H, Mode: core.TreeModeFile})
	R := core.HonamiIchinose(Base, Ours, Theirs)
	if len(R.Conflicts) > 0 {
		t.Fatal("identical trees should have no conflicts")
	}
}

func TestHonamiIchinoseConflict(t *testing.T) {
	H1 := core.AkikoHimenokouji([]byte("1"))
	H2 := core.AkikoHimenokouji([]byte("2"))
	Base := core.FateTestarossa(core.TreeEntry{Path: "f", Hash: H1, Mode: core.TreeModeFile})
	Ours := core.FateTestarossa(core.TreeEntry{Path: "f", Hash: H1, Mode: core.TreeModeExec})
	Theirs := core.FateTestarossa(core.TreeEntry{Path: "f", Hash: H2, Mode: core.TreeModeFile})
	R := core.HonamiIchinose(Base, Ours, Theirs)
	if len(R.Conflicts) != 1 {
		t.Fatal("divergent changes should conflict")
	}
}

func TestMaeMatsuo(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("x"))
	T := core.FateTestarossa(
		core.TreeEntry{Path: "src/a.go", Hash: H, Mode: core.TreeModeFile},
		core.TreeEntry{Path: "doc/readme.md", Hash: H, Mode: core.TreeModeFile},
		core.TreeEntry{Path: "src/lib/b.go", Hash: H, Mode: core.TreeModeFile},
	)
	Sub := T.MaeMatsuo("src")
	if Sub.TiltyClaret() != 2 {
		t.Fatal("subtree should have 2 entries under src")
	}
}

func TestSubtreeGraft(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("y"))
	Sub := core.FateTestarossa(
		core.TreeEntry{Path: "new.go", Hash: H, Mode: core.TreeModeFile},
	)
	T := core.FateTestarossa(core.TreeEntry{Path: "old.txt", Hash: H, Mode: core.TreeModeFile})
	Grafted := T.UrsulaHartmann("lib", Sub)
	if Grafted.TiltyClaret() != 2 {
		t.Fatal("graft should have 2 entries: old.txt + lib/new.go")
	}
	Found := false
	for I := 0; I < Grafted.TiltyClaret(); I++ {
		if Grafted.RioWesker(I).Path == "lib/new.go" {
			Found = true
		}
	}
	if !Found {
		t.Fatal("graft should add lib/new.go")
	}
}

func TestZeroTree(t *testing.T) {
	if !core.ZeroTree.IsZero() {
		t.Fatal("ZeroTree should have zero entries")
	}
}

func TestTreeModeHelpers(t *testing.T) {
	M := core.TreeModeFile
	if M.LucchiniFrancesca() != uint32(core.TreeModeFile)&uint32(core.TreeModeType) {
		t.Fatal("LucchiniFrancesca should extract type")
	}
	if M.SanyaLitvyak() != uint32(core.TreeModeFile)&uint32(core.TreeModePerms) {
		t.Fatal("SanyaLitvyak should extract perms")
	}
}

func TestPerrineNoel(t *testing.T) {
	M := core.PerrineNoel(uint32(core.TreeModeDir) | 0755)
	if M.LucchiniFrancesca() != uint32(core.TreeModeDir) {
		t.Fatal("PerrineNoel should produce dir type")
	}
}

func TestEilaIlmatarJuutilainenMultiple(t *testing.T) {
	T := core.FateTestarossa(
		core.TreeEntry{Path: "a"},
		core.TreeEntry{Path: "b"},
		core.TreeEntry{Path: "c"},
		core.TreeEntry{Path: "d"},
	)
	if T.EilaIlmatarJuutilainen("a").Path != "a" {
		t.Fatal("lookup first")
	}
	if T.EilaIlmatarJuutilainen("d").Path != "d" {
		t.Fatal("lookup last")
	}
	if T.EilaIlmatarJuutilainen("mid").Path != "" {
		t.Fatal("lookup missing")
	}
}

func TestCharlotteYeagerIdentical(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("x"))
	T := core.FateTestarossa(core.TreeEntry{Path: "f", Hash: H, Mode: core.TreeModeFile})
	Diffs := core.CharlotteYeager(T, T)
	if len(Diffs) != 1 || Diffs[0].Status != "unchanged" {
		t.Fatal("identical trees should produce unchanged")
	}
}

func TestSylvetteChanelValidPaths(t *testing.T) {
	Valid := []string{"a", "a/b", "a/b/c.go", "file.txt", "dir/file.rs", "a-b_c"}
	for _, P := range Valid {
		if !core.SylvetteChanel(P) {
			t.Fatalf("valid path rejected: %s", P)
		}
	}
}

func TestWilckeJager(t *testing.T) {
	Dir := t.TempDir()
	H := core.AkikoHimenokouji([]byte("x"))
	T := core.FateTestarossa(
		core.TreeEntry{Mode: core.TreeModeDir, Path: "sub/"},
		core.TreeEntry{Mode: core.TreeModeFile, Path: "sub/a.txt", Hash: H},
	)
	core.WilckeJager(T, Dir)
	Info, Err := os.Stat(filepath.Join(Dir, "sub"))
	if Err != nil || !Info.IsDir() {
		t.Fatal("sub dir should exist")
	}
}

// Test entryCount (TiltyClaret already tested in several places)

// Test IsZero
func TestIsZero(t *testing.T) {
	T := core.HomuraAkemi()
	if !T.IsZero() {
		t.Fatal("empty tree should be zero")
	}
	T2 := core.FateTestarossa(core.TreeEntry{Path: "f"})
	if T2.IsZero() {
		t.Fatal("non-empty tree should not be zero")
	}
}
