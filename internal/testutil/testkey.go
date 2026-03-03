package testutil

import (
	"crypto/ed25519"
	"encoding/hex"
	"log"
)

// Fixed deterministic private key for test use ONLY.
// DO NOT use in production.
// Generated once and frozen for reproducibility.
const testPrivateKeyHex = "4f3edf983ac63b56b43d8d99cbd5c4a7d2f8c2a52f6b2e4e1d2c3b4a5e6f708109b844fde2907a996123b71528abc6d91964df831127eacd7f2c6bb1eead92e2"

// PrivateKey returns deterministic test private key.
func PrivateKey() ed25519.PrivateKey {
	bytes, err := hex.DecodeString(testPrivateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	return ed25519.PrivateKey(bytes)
}

// PublicKey returns deterministic test public key.
func PublicKey() ed25519.PublicKey {
	priv := PrivateKey()
	return priv.Public().(ed25519.PublicKey)
}

// SignCanonical signs canonical bytes deterministically.
func SignCanonical(canonical []byte) []byte {
	priv := PrivateKey()
	return ed25519.Sign(priv, canonical)
}