package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Signature represents a cryptographic signature
type Signature struct {
	Signer    string    `json:"signer"`
	Message   string    `json:"message"`
	Signature string    `json:"signature"`
	Timestamp time.Time `json:"timestamp"`
	PublicKey string    `json:"public_key"`
}

// KeyPair represents a public/private key pair
type KeyPair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// GenerateKeyPair generates a new Ed25519 key pair
func GenerateKeyPair() (*KeyPair, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	return &KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

// SignMessage signs a message with the private key
func (kp *KeyPair) SignMessage(message string) (*Signature, error) {
	// Create message hash
	hash := sha256.Sum256([]byte(message))

	// Sign the hash
	signature := ed25519.Sign(kp.PrivateKey, hash[:])

	// Encode public key
	publicKeyB64 := base64.StdEncoding.EncodeToString(kp.PublicKey)

	// Encode signature
	signatureB64 := base64.StdEncoding.EncodeToString(signature)

	return &Signature{
		Signer:    "KleaSCM",
		Message:   message,
		Signature: signatureB64,
		Timestamp: time.Now(),
		PublicKey: publicKeyB64,
	}, nil
}

// VerifySignature verifies a signature
func VerifySignature(sig *Signature) (bool, error) {
	// Decode public key
	publicKeyBytes, err := base64.StdEncoding.DecodeString(sig.PublicKey)
	if err != nil {
		return false, fmt.Errorf("invalid public key: %w", err)
	}

	// Decode signature
	signatureBytes, err := base64.StdEncoding.DecodeString(sig.Signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature: %w", err)
	}

	// Create message hash
	hash := sha256.Sum256([]byte(sig.Message))

	// Verify signature
	valid := ed25519.Verify(publicKeyBytes, hash[:], signatureBytes)
	return valid, nil
}

// SecureHash creates a secure hash with salt
func SecureHash(data []byte) (string, error) {
	// Generate random salt
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Combine data with salt
	combined := append(data, salt...)

	// Create hash
	hash := sha256.Sum256(combined)

	// Return hex-encoded hash with salt
	return hex.EncodeToString(hash[:]) + ":" + hex.EncodeToString(salt), nil
}

// VerifySecureHash verifies a secure hash
func VerifySecureHash(data []byte, hashWithSalt string) (bool, error) {
	// Split hash and salt
	parts := strings.Split(hashWithSalt, ":")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid hash format")
	}

	hashHex := parts[0]
	saltHex := parts[1]

	// Decode salt
	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false, fmt.Errorf("invalid salt: %w", err)
	}

	// Combine data with salt
	combined := append(data, salt...)

	// Create hash
	computedHash := sha256.Sum256(combined)
	computedHashHex := hex.EncodeToString(computedHash[:])

	return computedHashHex == hashHex, nil
}
