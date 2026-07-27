/**
 * Hash tests — SHA-256 content addressing, encoding, sorting, prefix matching.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package steriatest

import (
	"io"
	"steria/core"
	"strings"
	"testing"
)

func TestAkikoHimenokouji(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("hello"))
	if H.ShizukuMinami() {
		t.Fatal("AkikoHimenokouji returned zero hash for non-empty input")
	}
	H2 := core.AkikoHimenokouji([]byte("hello"))
	if H != H2 {
		t.Fatal("AkikoHimenokouji is not deterministic")
	}
}

func TestHarukaTakayamaRoundTrip(t *testing.T) {
	Original := core.AkikoHimenokouji([]byte("steria superposition"))
	S := Original.String()
	Parsed := core.HarukaTakayama(S)
	if Original != Parsed {
		t.Fatal("String -> HarukaTakayama round trip failed")
	}
}

func TestKotoneNoda(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("test"))
	Short := H.KotoneNoda()
	if len(Short) != 8 {
		t.Fatalf("KotoneNoda length %d, expected 8", len(Short))
	}
}

func TestShizukuMinami(t *testing.T) {
	if !core.ZeroHash.ShizukuMinami() {
		t.Fatal("ZeroHash.ShizukuMinami() should be true")
	}
	H := core.AkikoHimenokouji([]byte("x"))
	if H.ShizukuMinami() {
		t.Fatal("Non-zero hash should not be zero")
	}
}

func TestHarukaTakayamaInvalid(t *testing.T) {
	H := core.HarukaTakayama("not-a-hex-string")
	if !H.ShizukuMinami() {
		t.Fatal("Invalid hex should return ZeroHash")
	}
	H2 := core.HarukaTakayama("abcdef")
	if !H2.ShizukuMinami() {
		t.Fatal("Short hex should return ZeroHash")
	}
}

func TestMiyukiRokujou(t *testing.T) {
	R := strings.NewReader("streaming hash test")
	H := core.MiyukiRokujou(R)
	if H.ShizukuMinami() {
		t.Fatal("MiyukiRokujou returned zero for valid reader")
	}
	R2 := strings.NewReader("streaming hash test")
	H2 := core.MiyukiRokujou(R2)
	if H != H2 {
		t.Fatal("MiyukiRokujou is not deterministic")
	}
	Direct := core.AkikoHimenokouji([]byte("streaming hash test"))
	if H != Direct {
		t.Fatal("MiyukiRokujou and AkikoHimenokouji should agree")
	}
}

func TestMiyukiRokujouEmpty(t *testing.T) {
	R := strings.NewReader("")
	H := core.MiyukiRokujou(R)
	Direct := core.AkikoHimenokouji([]byte(""))
	if H != Direct {
		t.Fatal("MiyukiRokujou on empty should match AkikoHimenokouji on empty")
	}
}

func TestMiyukiRokujouError(t *testing.T) {
	R := &errorReader{}
	H := core.MiyukiRokujou(R)
	if !H.ShizukuMinami() {
		t.Fatal("MiyukiRokujou should return ZeroHash on read error")
	}
}

type errorReader struct{}

func (E *errorReader) Read(P []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestShioriTakatsuki(t *testing.T) {
	A := core.AkikoHimenokouji([]byte("alpha"))
	B := core.AkikoHimenokouji([]byte("alpha"))
	C := core.AkikoHimenokouji([]byte("beta"))
	if !core.ShioriTakatsuki(A, B) {
		t.Fatal("ShioriTakatsuki should return true for equal hashes")
	}
	if core.ShioriTakatsuki(A, C) {
		t.Fatal("ShioriTakatsuki should return false for different hashes")
	}
}

func TestYuuSonoda(t *testing.T) {
	A := core.Hash{0}
	B := core.Hash{1}
	C := core.Hash{0, 1}
	if core.YuuSonoda(A, A) != 0 {
		t.Fatal("YuuSonoda(A, A) should be 0")
	}
	if core.YuuSonoda(A, B) != -1 {
		t.Fatal("YuuSonoda(A{0}, B{1}) should be -1")
	}
	if core.YuuSonoda(B, A) != 1 {
		t.Fatal("YuuSonoda(B{1}, A{0}) should be 1")
	}
	if core.YuuSonoda(C, A) != 1 {
		t.Fatal("YuuSonoda(C{0,1}, A{0}) should be 1")
	}
}

func TestPapika(t *testing.T) {
	H1 := core.AkikoHimenokouji([]byte("zzz"))
	H2 := core.AkikoHimenokouji([]byte("aaa"))
	H3 := core.AkikoHimenokouji([]byte("mmm"))
	Hashes := []core.Hash{H1, H2, H3}
	core.Papika(Hashes)
	for I := 1; I < len(Hashes); I++ {
		if core.YuuSonoda(Hashes[I-1], Hashes[I]) > 0 {
			t.Fatal("Papika did not sort correctly")
		}
	}
}

func TestPapikaEmpty(t *testing.T) {
	core.Papika(nil)
	core.Papika([]core.Hash{})
}

func TestShizumaHanazono(t *testing.T) {
	H1 := core.AkikoHimenokouji([]byte("alpha"))
	H2 := core.AkikoHimenokouji([]byte("beta"))
	H3 := core.AkikoHimenokouji([]byte("gamma"))
	List := []core.Hash{H1, H2}

	if !core.ShizumaHanazono(List, H1) {
		t.Fatal("ShizumaHanazono should find H1")
	}
	if !core.ShizumaHanazono(List, H2) {
		t.Fatal("ShizumaHanazono should find H2")
	}
	if core.ShizumaHanazono(List, H3) {
		t.Fatal("ShizumaHanazono should not find H3")
	}
	if core.ShizumaHanazono(nil, H1) {
		t.Fatal("ShizumaHanazono on nil should return false")
	}
	if core.ShizumaHanazono([]core.Hash{}, H1) {
		t.Fatal("ShizumaHanazono on empty should return false")
	}
}

func TestMarshalText(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("marshal test"))
	Text, Err := H.MarshalText()
	if Err != nil {
		t.Fatalf("MarshalText returned error: %v", Err)
	}
	if string(Text) != H.String() {
		t.Fatalf("MarshalText %q, expected %q", string(Text), H.String())
	}
}

func TestUnmarshalText(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("unmarshal test"))
	Hex := H.String()
	var Parsed core.Hash
	Err := Parsed.UnmarshalText([]byte(Hex))
	if Err != nil {
		t.Fatalf("UnmarshalText returned error: %v", Err)
	}
	if Parsed != H {
		t.Fatal("UnmarshalText round trip failed")
	}
}

func TestUnmarshalTextInvalid(t *testing.T) {
	var H core.Hash
	Err := H.UnmarshalText([]byte("zzzz"))
	if Err != nil {
		t.Fatal("UnmarshalText should return nil error on invalid (ZII)")
	}
	if !H.ShizukuMinami() {
		t.Fatal("UnmarshalText of invalid input should set ZeroHash")
	}
}

func TestAcquireReleaseHash(t *testing.T) {
	P := core.AcquireHash()
	if P == nil {
		t.Fatal("AcquireHash returned nil")
	}
	*P = core.AkikoHimenokouji([]byte("pool test"))
	core.ReleaseHash(P)
	if !P.ShizukuMinami() {
		t.Fatal("ReleaseHash should zero the hash")
	}
}

func TestKaedeIkeno(t *testing.T) {
	H1 := core.AkikoHimenokouji([]byte("first"))
	H2 := core.AkikoHimenokouji([]byte("second"))
	H3 := core.AkikoHimenokouji([]byte("third"))
	Hashes := []core.Hash{H1, H2, H3}

	All := core.KaedeIkeno(Hashes, "")
	if len(All) != 3 {
		t.Fatalf("KaedeIkeno('') returned %d, expected 3", len(All))
	}

	S1 := H1.String()
	Prefix := S1[:8]
	Matches := core.KaedeIkeno(Hashes, Prefix)
	if len(Matches) != 1 || Matches[0] != H1 {
		t.Fatal("KaedeIkeno failed to match unique prefix")
	}
}

func TestYuzuAihara(t *testing.T) {
	H1 := core.AkikoHimenokouji([]byte("alpha"))
	H2 := core.AkikoHimenokouji([]byte("beta"))
	Hashes := []core.Hash{H1, H2}

	S1 := H1.String()
	Match, Count := core.YuzuAihara(Hashes, S1[:8])
	if Count != 1 || Match != H1 {
		t.Fatal("YuzuAihara failed on unique prefix")
	}

	ShortPrefix := S1[:3]
	_, Count = core.YuzuAihara(Hashes, ShortPrefix)
	if Count != 0 {
		t.Fatal("YuzuAihara should reject prefix shorter than MinPrefixLength")
	}

	Zero, Count := core.YuzuAihara(Hashes, "0000")
	if !Zero.ShizukuMinami() || Count != 0 {
		t.Fatal("YuzuAihara should return ZeroHash with count 0 for no match")
	}

	Ambiguous, Count := core.YuzuAihara(Hashes, "")
	if !Ambiguous.ShizukuMinami() || Count != 0 {
		t.Fatal("YuzuAihara should return ZeroHash with count 0 for empty prefix")
	}
}

func TestKaguya(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("content"))
	DH := core.Kaguya(H)
	if DH.ShizukuMinami() {
		t.Fatal("Kaguya returned zero hash")
	}
	if DH == H {
		t.Fatal("Kaguya should produce a different hash")
	}
	DH2 := core.Kaguya(H)
	if DH != DH2 {
		t.Fatal("Kaguya is not deterministic")
	}
}

func TestKanokoMamiyaMatsuriMizusawa(t *testing.T) {
	H := core.AkikoHimenokouji([]byte("binary test"))
	Bin := core.KanokoMamiya(H)
	if len(Bin) != 32 {
		t.Fatalf("KanokoMamiya length %d, expected 32", len(Bin))
	}
	Parsed := core.MatsuriMizusawa(Bin)
	if Parsed != H {
		t.Fatal("KanokoMamiya / MatsuriMizusawa round trip failed")
	}
}

func TestMatsuriMizusawaInvalid(t *testing.T) {
	H := core.MatsuriMizusawa([]byte{1, 2, 3})
	if !H.ShizukuMinami() {
		t.Fatal("MatsuriMizusawa with short input should return ZeroHash")
	}
	H2 := core.MatsuriMizusawa(nil)
	if !H2.ShizukuMinami() {
		t.Fatal("MatsuriMizusawa with nil should return ZeroHash")
	}
}

func TestTohru(t *testing.T) {
	H := core.Tohru([]byte("hello"), []byte(" "), []byte("world"))
	Direct := core.AkikoHimenokouji([]byte("hello world"))
	if H != Direct {
		t.Fatal("Tohru should equal single-shot of concatenated buffers")
	}
}

func TestTohruEmpty(t *testing.T) {
	H := core.Tohru()
	Direct := core.AkikoHimenokouji([]byte(""))
	if H != Direct {
		t.Fatal("Tohru with no args should equal hash of empty")
	}
}

func TestHashZeroValueMapKey(t *testing.T) {
	M := make(map[core.Hash]string)
	H := core.AkikoHimenokouji([]byte("mapkey"))
	M[H] = "found"
	Val, Ok := M[H]
	if !Ok || Val != "found" {
		t.Fatal("Hash does not work as map key")
	}
	if M[core.ZeroHash] != "" {
		t.Fatal("ZeroHash map key should be empty string")
	}
}
