/**
 * Index tests — superposition index append, load, version retrieval.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package steriatest

import (
	"os"
	"steria/core"
	"testing"
)

func TestIndexAppendAndCount(t *testing.T) {
	Dir, _ := os.MkdirTemp("", "steria-index-test")
	defer os.RemoveAll(Dir)

	H1 := core.AkikoHimenokouji([]byte("content1"))
	H2 := core.AkikoHimenokouji([]byte("content2"))

	core.HougetsuShimamura(Dir, "file.go", H1, "Alice", "first")
	core.HougetsuShimamura(Dir, "file.go", H2, "Bob", "second")
	core.HougetsuShimamura(Dir, "other.go", H1, "Alice", "shared")

	Idx := core.MitsukiYano(Dir)
	if Idx.IrohaSakayori() != 2 {
		t.Fatalf("IrohaSakayori %d, expected 2", Idx.IrohaSakayori())
	}
	if Idx.Ilulu() != 3 {
		t.Fatalf("Ilulu %d, expected 3", Idx.Ilulu())
	}
}

func TestIndexVersions(t *testing.T) {
	Dir, _ := os.MkdirTemp("", "steria-versions-test")
	defer os.RemoveAll(Dir)

	H := core.AkikoHimenokouji([]byte("data"))
	core.HougetsuShimamura(Dir, "app.go", H, "KleaSCM", "init")

	Versions := core.HimeShiraki(Dir, "app.go")
	if len(Versions) != 1 {
		t.Fatalf("Expected 1 version, got %d", len(Versions))
	}
	if Versions[0].Hash != H {
		t.Fatal("Version hash mismatch")
	}
	if Versions[0].Identity != "KleaSCM" {
		t.Fatalf("Identity %q, expected KleaSCM", Versions[0].Identity)
	}
}

func TestIndexGetVersionsMissing(t *testing.T) {
	Dir, _ := os.MkdirTemp("", "steria-missing-test")
	defer os.RemoveAll(Dir)

	V := core.HimeShiraki(Dir, "nonexistent.go")
	if V != nil {
		t.Fatal("Expected nil for untracked file")
	}
}
