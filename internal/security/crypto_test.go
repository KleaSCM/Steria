package security

import (
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}
	if kp.PublicKey == nil || kp.PrivateKey == nil {
		t.Fatalf("Key pair fields are nil")
	}
}

func TestSignAndVerifyMessage(t *testing.T) {
	kp, _ := GenerateKeyPair()
	sig, err := kp.SignMessage("hello")
	if err != nil {
		t.Fatalf("SignMessage failed: %v", err)
	}
	valid, err := VerifySignature(sig)
	if err != nil {
		t.Fatalf("VerifySignature failed: %v", err)
	}
	if !valid {
		t.Errorf("Signature should be valid")
	}
}

func TestVerifySignature_Invalid(t *testing.T) {
	kp, _ := GenerateKeyPair()
	sig, _ := kp.SignMessage("hello")
	sig.Signature = "invalid"
	valid, err := VerifySignature(sig)
	if err == nil || valid {
		t.Errorf("Expected invalid signature error")
	}
}

func TestSecureHashAndVerify(t *testing.T) {
	data := []byte("test")
	hash, err := SecureHash(data)
	if err != nil {
		t.Fatalf("SecureHash failed: %v", err)
	}
	ok, err := VerifySecureHash(data, hash)
	if err != nil {
		t.Fatalf("VerifySecureHash failed: %v", err)
	}
	if !ok {
		t.Errorf("Secure hash should verify")
	}
}

func TestVerifySecureHash_Invalid(t *testing.T) {
	ok, err := VerifySecureHash([]byte("test"), "invalid")
	if err == nil || ok {
		t.Errorf("Expected invalid hash error")
	}
}

func BenchmarkGenerateKeyPair(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateKeyPair()
	}
}
