package signature

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"
)

func TestSignatureVerification(t *testing.T) {

	message := []byte("dip-canonical-test")

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	signature := ed25519.Sign(privateKey, message)

	pubHex := hex.EncodeToString(publicKey)
	sigHex := hex.EncodeToString(signature)

	valid, err := VerifySignature(pubHex, sigHex, message)
	if err != nil {
		t.Fatalf("Unexpected error during verification: %v", err)
	}

	if !valid {
		t.Fatalf("Expected valid signature but verification failed")
	}
}

func TestInvalidSignature(t *testing.T) {

	message := []byte("dip-canonical-test")

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	signature := ed25519.Sign(privateKey, message)

	signature[0] ^= 0xFF

	pubHex := hex.EncodeToString(publicKey)
	sigHex := hex.EncodeToString(signature)

	valid, err := VerifySignature(pubHex, sigHex, message)
	if err != nil {
		t.Fatalf("Unexpected error during verification: %v", err)
	}

	if valid {
		t.Fatalf("Expected invalid signature but verification passed")
	}
}