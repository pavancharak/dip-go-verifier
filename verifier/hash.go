package verifier

import (
	"crypto/sha256"
	"encoding/hex"
)

func ComputeArtifactID(canonical []byte) string {

	hash := sha256.Sum256(canonical)
	return hex.EncodeToString(hash[:])
}