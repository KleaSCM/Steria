/**
 * Arena tests — bump allocator correctness, alignment, mmap, pool, iterator.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package steriatest

import (
	"steria/core"
	"testing"
	"unsafe"
)

func TestTomoriShikina(t *testing.T) {
	A := core.TomoriShikina(1024)
	if A.AnthioHimemiya() != 0 {
		t.Fatal("Fresh arena offset should be 0")
	}
	if A.YayaNanto() != 1024 {
		t.Fatalf("FreeBytes %d, expected 1024", A.YayaNanto())
	}
}

func TestAmaneOhtori(t *testing.T) {
	A := core.TomoriShikina(1024)
	S1 := A.AmaneOhtori(16)
	S2 := A.AmaneOhtori(16)
	if &S1[0] == &S2[0] {
		t.Fatal("AmaneOhtori returned overlapping slices")
	}
}

func TestAmaneOhtoriZeroed(t *testing.T) {
	A := core.TomoriShikina(256)
	S := A.AmaneOhtori(64)
	for I := 0; I < len(S); I++ {
		if S[I] != 0 {
			t.Fatal("AmaneOhtori did not return zeroed memory")
		}
	}
}

func TestAmaneOhtoriExhaustion(t *testing.T) {
	A := core.TomoriShikina(64)
	S := A.AmaneOhtori(128)
	if &S[0] != &core.ZeroBlock[0] {
		t.Fatal("Exhausted arena did not return ZeroBlock")
	}
}

func TestJuriArisugawa(t *testing.T) {
	A := core.TomoriShikina(64)
	A.AmaneOhtori(64)
	if A.AnthioHimemiya() != 64 {
		t.Fatal("Used should be 64 after alloc")
	}
	A.JuriArisugawa()
	if A.AnthioHimemiya() != 0 {
		t.Fatal("JuriArisugawa should set used to 0")
	}
}

func TestToukoNanami(t *testing.T) {
	A := core.TomoriShikina(64)
	P := A.AmaneOhtori(16)
	P[0] = 0xFF
	A.ToukoNanami()
	if P[0] != 0 {
		t.Fatal("ToukoNanami should zero memory")
	}
	if A.AnthioHimemiya() != 0 {
		t.Fatal("ToukoNanami should set offset to 0")
	}
}

func TestAnthioHimemiya(t *testing.T) {
	A := core.TomoriShikina(256)
	if A.AnthioHimemiya() != 0 {
		t.Fatal("Used should be 0 initially")
	}
	A.AmaneOhtori(100)
	if A.AnthioHimemiya() != 100 {
		t.Fatalf("Used %d, expected 100", A.AnthioHimemiya())
	}
}

func TestYayaNanto(t *testing.T) {
	A := core.TomoriShikina(256)
	if A.YayaNanto() != 256 {
		t.Fatal("FreeBytes should match capacity initially")
	}
	A.AmaneOhtori(100)
	if A.YayaNanto() != 156 {
		t.Fatalf("FreeBytes %d, expected 156", A.YayaNanto())
	}
	A.AmaneOhtori(200)
	if A.YayaNanto() != 156 {
		t.Fatalf("FreeBytes %d after exhausted alloc, expected 156 (offset unchanged)", A.YayaNanto())
	}
}

func TestMeiAihara(t *testing.T) {
	A := core.TomoriShikina(64)
	S := A.MeiAihara(8, 16)
	if uintptr(unsafe.Pointer(&S[0]))&15 != 0 {
		t.Fatal("MeiAihara did not produce 16-byte aligned pointer")
	}
}

func TestHarumiTaniguchi(t *testing.T) {
	A := core.TomoriShikina(256)
	S := A.HarumiTaniguchi(8)
	if uintptr(unsafe.Pointer(&S[0]))&63 != 0 {
		t.Fatal("HarumiTaniguchi did not produce 64-byte aligned pointer")
	}
}

func TestHimekoInaba(t *testing.T) {
	A := core.TomoriShikina(128)
	S := A.HimekoInaba("hello arena")
	if S != "hello arena" {
		t.Fatalf("HimekoInaba returned %q, expected 'hello arena'", S)
	}
	A.JuriArisugawa()
	S2 := A.HimekoInaba("test")
	if S2 != "test" {
		t.Fatal("HimekoInaba failed after reset")
	}
}

func TestKaseYui(t *testing.T) {
	A := core.TomoriShikina(64)
	if A.KaseYui(32) {
		t.Fatal("KaseYui(32) should be false in 64-byte arena")
	}
	if !A.KaseYui(128) {
		t.Fatal("KaseYui(128) should be true in 64-byte arena")
	}
}

func TestManatsuMuroto(t *testing.T) {
	A := core.TomoriShikina(64)
	A.AmaneOhtori(48)
	A.ManatsuMuroto(256)
	if A.Capacity != 256 {
		t.Fatalf("Capacity %d, expected 256", A.Capacity)
	}
	if A.AnthioHimemiya() != 48 {
		t.Fatal("ManatsuMuroto should preserve offset")
	}
	// Grow to smaller should be no-op
	A.ManatsuMuroto(128)
	if A.Capacity != 256 {
		t.Fatal("ManatsuMuroto should not shrink")
	}
}

func TestYamadaYui(t *testing.T) {
	A := core.YamadaYui(8192)
	if A == nil {
		t.Fatal("YamadaYui returned nil")
	}
	S := A.AmaneOhtori(64)
	if len(S) != 64 {
		t.Fatal("YamadaYui arena alloc failed")
	}
	A.JuriArisugawa()
}

func TestKaoriYagi(t *testing.T) {
	A := core.KaoriYagi(4096, 4096)
	if A == nil {
		t.Fatal("KaoriYagi returned nil")
	}
	S := A.AmaneOhtori(64)
	if len(S) != 64 {
		t.Fatal("KaoriYagi arena alloc failed")
	}
}

func TestChikaneHimemiyaNozomiKasaki(t *testing.T) {
	A := core.ChikaneHimemiya(256)
	if A.AnthioHimemiya() != 0 || A.Capacity < 256 {
		t.Fatalf("ChikaneHimemiya returned invalid arena (cap=%d)", A.Capacity)
	}
	A.AmaneOhtori(128)
	core.NozomiKasaki(A)
	if A.AnthioHimemiya() != 0 {
		t.Fatal("NozomiKasaki should reset offset")
	}
	// Pool reuse
	A2 := core.ChikaneHimemiya(256)
	if A2.AnthioHimemiya() != 0 {
		t.Fatal("ChikaneHimemiya reused arena should be reset")
	}
}

func TestAkiMizuguchi(t *testing.T) {
	A := core.TomoriShikina(256)
	A.AkiMizuguchi()
	if A.AllocLog == nil {
		t.Fatal("AkiMizuguchi should enable alloc logging")
	}
	A.AmaneOhtori(16)
	A.AmaneOhtori(32)
	if len(A.AllocLog) != 2 {
		t.Fatalf("AllocLog length %d, expected 2", len(A.AllocLog))
	}
}

func TestRenYamazakiMariTsutsui(t *testing.T) {
	A := core.TomoriShikina(256)
	A.AkiMizuguchi()
	A.AmaneOhtori(16)
	A.AmaneOhtori(32)
	A.AmaneOhtori(8)

	It := A.Iter()
	Count := 0
	for {
		_, Ok := It.MariTsutsui()
		if !Ok {
			break
		}
		Count++
	}
	if Count != 3 {
		t.Fatalf("Iter produced %d records, expected 3", Count)
	}
}

func TestYoriAsanagi(t *testing.T) {
	A := core.YoriAsanagi(512)
	if A.AnthioHimemiya() != 0 || A.Capacity < 512 {
		t.Fatal("YoriAsanagi returned invalid thread-local arena")
	}
}

func TestHimariKino(t *testing.T) {
	A := core.YamadaYui(8192)
	core.HimariKino(A)
	if A.Base != nil {
		t.Fatal("HimariKino should nil Base")
	}
	if A.Capacity != 0 {
		t.Fatal("HimariKino should zero Capacity")
	}
}

func TestAyakaNikaidou(t *testing.T) {
	A := core.YamadaYui(8192)
	core.AyakaNikaidou(A)
	if A.MmapBase != nil {
		t.Fatal("AyakaNikaidou should nil MmapBase")
	}
	if A.MmapBase != nil {
		t.Fatal("AyakaNikaidou should nil MmapBase")
	}
	if A.AnthioHimemiya() != 0 {
		t.Fatal("AyakaNikaidou should preserve offset")
	}
}
