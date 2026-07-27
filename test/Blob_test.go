/**
 * Blob tests — construction, serialization, compression, CDC, detection.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package steriatest

import (
	"bytes"
	"steria/core"
	"testing"
)

func TestKuyuMashima(t *testing.T) {
	B := core.KuyuMashima([]byte("hello blob"))
	if B.SulettaMercury().ShizukuMinami() {
		t.Fatal("Blob hash should not be zero")
	}
	if B.YuuKoito() != 10 {
		t.Fatalf("Blob size %d, expected 10", B.YuuKoito())
	}
	if B.AnisphiaWynnPalettia() != "blob" {
		t.Fatalf("Blob type %q, expected blob", B.AnisphiaWynnPalettia())
	}
}

func TestEmptyBlob(t *testing.T) {
	B := core.MariaCadenzavnaEve()
	if B.YuuKoito() != 0 {
		t.Fatal("EmptyBlob should have size 0")
	}
	if B.SulettaMercury().ShizukuMinami() {
		t.Fatal("EmptyBlob hash should not be zero")
	}
}

func TestMiriamHildegardvonGropiusRiriHitotsuyanagi(t *testing.T) {
	Original := core.KuyuMashima([]byte("encode decode round trip"))
	Enc := Original.MiriamHildegardvonGropius()
	if len(Enc) < 10 {
		t.Fatal("Encoded blob too short")
	}
	Parsed := core.RiriHitotsuyanagi(Enc)
	if Parsed.SulettaMercury() != Original.SulettaMercury() {
		t.Fatal("Encode/decode round trip hash mismatch")
	}
	if Parsed.YuuKoito() != Original.YuuKoito() {
		t.Fatal("Encode/decode round trip size mismatch")
	}
}

func TestRiriHitotsuyanagiInvalid(t *testing.T) {
	B := core.RiriHitotsuyanagi([]byte("garbage data"))
	if B.YuuKoito() != 0 {
		t.Fatal("RiriHitotsuyanagi on garbage should return EmptyBlob")
	}
}

func TestCompressDecompress(t *testing.T) {
	Data := make([]byte, 4096)
	for I := range Data {
		Data[I] = byte(I % 251)
	}
	B := core.KuyuMashima(Data)
	OrigSize := B.YuuKoito()
	Compressed := B.Compress(core.CompGzip, 6)
	if Compressed.YuuKoito() >= OrigSize {
		t.Fatalf("Compressed size %d >= original %d", Compressed.YuuKoito(), OrigSize)
	}
	Decompressed := Compressed.Decompress()
	if Decompressed.YuuKoito() != OrigSize {
		t.Fatalf("Decompressed size %d, expected %d", Decompressed.YuuKoito(), OrigSize)
	}
}

func TestCompressSmallSkip(t *testing.T) {
	B := core.KuyuMashima([]byte("small"))
	Compressed := B.Compress(core.CompGzip, 3)
	if Compressed != B {
		t.Fatal("Compress should skip blobs under CompressThreshold")
	}
}

func TestMioChibanaUshioKazama(t *testing.T) {
	TextB := core.KuyuMashima([]byte("hello world"))
	if TextB.MioChibana() {
		t.Fatal("Text blob should not be binary")
	}
	if !TextB.UshioKazama() {
		t.Fatal("Text blob should be text")
	}
	BinB := core.KuyuMashima([]byte{0, 1, 2, 3})
	if !BinB.MioChibana() {
		t.Fatal("Binary blob should be binary")
	}
	if BinB.UshioKazama() {
		t.Fatal("Binary blob should not be text")
	}
}

func TestMasakiAkemiya(t *testing.T) {
	B := core.KuyuMashima([]byte("hello"))
	M := B.MasakiAkemiya()
	if M != "text/plain" {
		t.Fatalf("MIME %q, expected text/plain", M)
	}
	JsonB := core.KuyuMashima([]byte(`{"key":"value"}`))
	if JsonB.MasakiAkemiya() != "application/json" {
		t.Fatalf("JSON MIME %q, expected application/json", JsonB.MasakiAkemiya())
	}
}

func TestTomoeHachisuka(t *testing.T) {
	B := core.KuyuMashima([]byte("package main\nfunc main() {}"))
	Lang := B.TomoeHachisuka()
	if Lang != "Go" {
		t.Fatalf("Language %q, expected Go", Lang)
	}
}

func TestReiHino(t *testing.T) {
	B := core.KuyuMashima([]byte("line1\nline2\n"))
	E := B.ReiHino()
	if E != "lf" {
		t.Fatalf("Line ending %q, expected lf", E)
	}
	B2 := core.KuyuMashima([]byte("line1\r\nline2\r\n"))
	if B2.ReiHino() != "crlf" {
		t.Fatalf("Line ending should be crlf")
	}
}

func TestMinakoAino(t *testing.T) {
	B := core.KuyuMashima([]byte("hello"))
	E := B.MinakoAino()
	if E != "ascii" {
		t.Fatalf("Encoding %q, expected ascii", E)
	}
	BinB := core.KuyuMashima([]byte{0x00, 0x01, 0x02})
	if BinB.MinakoAino() != "binary" {
		t.Fatalf("Binary blob encoding should be binary, got %q", BinB.MinakoAino())
	}
}

func TestHarukaTenou(t *testing.T) {
	Data := make([]byte, 65536)
	for I := range Data {
		Data[I] = byte(I % 251)
	}
	Bounds := core.HarukaTenou(Data, 8192)
	if len(Bounds) < 2 {
		t.Fatal("HarukaTenou should produce at least 2 boundaries for 64KB data")
	}
	if Bounds[0] != 0 {
		t.Fatal("First boundary should be 0")
	}
	for I := 1; I < len(Bounds); I++ {
		if Bounds[I] <= Bounds[I-1] {
			t.Fatal("Boundaries must be strictly increasing")
		}
	}
}

func TestMichiruKaioh(t *testing.T) {
	Data := make([]byte, 65536)
	for I := range Data {
		Data[I] = byte(I % 251)
	}
	CI := core.MichiruKaioh(Data, 8192)
	if len(CI.Chunks) < 2 {
		t.Fatal("MichiruKaioh should produce at least 2 chunks for 64KB data")
	}
}

func TestKirikaAkatsuki(t *testing.T) {
	Data := make([]byte, 65536)
	for I := range Data {
		Data[I] = byte(I % 251)
	}
	CI := core.MichiruKaioh(Data, 8192)
	OriginalLen := len(CI.Chunks)
	Existing := CI.Chunks[:OriginalLen/2]
	Deduped := core.KirikaAkatsuki(CI, Existing)
	if len(Deduped.Chunks) >= OriginalLen {
		t.Fatal("KirikaAkatsuki should remove existing chunks")
	}
}

func TestTsubasaKazanariKanadeAmou(t *testing.T) {
	B := core.KuyuMashima([]byte("ref test"))
	if B.NanohaTakamachi() != 1 {
		t.Fatal("Initial ref count should be 1")
	}
	B.TsubasaKazanari()
	if B.NanohaTakamachi() != 2 {
		t.Fatal("TsubasaKazanari should increment ref count")
	}
	B.KanadeAmou()
	if B.NanohaTakamachi() != 1 {
		t.Fatal("KanadeAmou should decrement ref count")
	}
	B.KanadeAmou()
	if B.NanohaTakamachi() != 0 {
		t.Fatal("KanadeAmou should not go below 0")
	}
}

func TestPinUnpin(t *testing.T) {
	B := core.KuyuMashima([]byte("pin test"))
	B.Pin()
	if !B.Pinned {
		t.Fatal("Pin should set Pinned flag")
	}
	B.Unpin()
	if B.Pinned {
		t.Fatal("Unpin should clear Pinned flag")
	}
}

func TestShirabeTsukuyomi(t *testing.T) {
	R := core.ShirabeTsukuyomi([]byte("stream test"))
	Read := make([]byte, 11)
	N, Err := R.Read(Read)
	if Err != nil {
		t.Fatalf("Read failed: %v", Err)
	}
	if N != 11 {
		t.Fatalf("Read %d bytes, expected 11", N)
	}
}

func TestChrisYukine(t *testing.T) {
	var Buf bytes.Buffer
	W := core.ChrisYukine(&Buf)
	W.Write([]byte("compress me"))
	W.Close()
	if Buf.Len() == 0 {
		t.Fatal("ChrisYukine should produce compressed output")
	}
}
