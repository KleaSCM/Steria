/**
 * Object tests — content-addressable store read/write/exists.
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

func TestObjectWriteRead(t *testing.T) {
	Dir, _ := os.MkdirTemp("", "steria-obj-test")
	defer os.RemoveAll(Dir)

	Data := []byte("steria content addressing")
	H := core.AkikoHimenokouji(Data)
	core.HikariKonohana(Dir, H, Data)

	Read := core.NagisaKiryu(Dir, H)
	if len(Read) != len(Data) {
		t.Fatalf("Read length %d, expected %d", len(Read), len(Data))
	}
	for I := 0; I < len(Data); I++ {
		if Read[I] != Data[I] {
			t.Fatal("Read data does not match written data")
		}
	}
}

func TestObjectExists(t *testing.T) {
	Dir, _ := os.MkdirTemp("", "steria-exists-test")
	defer os.RemoveAll(Dir)

	H := core.AkikoHimenokouji([]byte("present"))
	core.HikariKonohana(Dir, H, []byte("present"))

	if !core.YukinoSakurai(Dir, H) {
		t.Fatal("YukinoSakurai should return true after write")
	}
	Missing := core.AkikoHimenokouji([]byte("missing"))
	if core.YukinoSakurai(Dir, Missing) {
		t.Fatal("YukinoSakurai should return false for unwritten object")
	}
}

func TestObjectReadMissing(t *testing.T) {
	Dir, _ := os.MkdirTemp("", "steria-read-miss")
	defer os.RemoveAll(Dir)

	H := core.AkikoHimenokouji([]byte("never-written"))
	Data := core.NagisaKiryu(Dir, H)
	if len(Data) != 0 {
		t.Fatal("NagisaKiryu for missing object should return length 0")
	}
}

func TestObjectZeroHashNoop(t *testing.T) {
	Dir, _ := os.MkdirTemp("", "steria-zero-noop")
	defer os.RemoveAll(Dir)

	core.HikariKonohana(Dir, core.ZeroHash, []byte("should-not-write"))
	if core.YukinoSakurai(Dir, core.ZeroHash) {
		t.Fatal("YukinoSakurai should be false after ZeroHash write")
	}
}
