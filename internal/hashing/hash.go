package hashing

import (
	"crypto/sha256"
	"encoding/hex"
)

// ComputeSHA256 returns the SHA-256 hash of the given byte slice in hex.
func ComputeSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}