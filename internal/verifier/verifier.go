package verifier

import (
	"encoding/json"
	"errors"

	"github.com/pavancharak/dip-go-verifier/internal/artifact"
	"github.com/pavancharak/dip-go-verifier/internal/canonicalization"
	"github.com/pavancharak/dip-go-verifier/internal/signature"
)

// VerifyArtifact verifies a DIP artifact using provided public keys.
// Returns true if verification succeeds, false otherwise.
func VerifyArtifact(data []byte, pubKeys map[string]string) (bool, error) {

	// Parse artifact
	da, err := artifact.ParseDecisionArtifact(data)
	if err != nil {
		return false, err
	}

	// Canonicalize
	canonical := canonicalization.Canonicalize(data)

	// Lookup public key
	pubKeyHex, ok := pubKeys[da.ArtifactID]
	if !ok {
		return false, errors.New("public key not found for artifact")
	}

	// Verify signature
	valid, err := signature.VerifySignature(pubKeyHex, da.Signature, canonical)
	if err != nil {
		return false, err
	}

	return valid, nil
}

// LoadPublicKeys loads public keys from JSON file bytes.
func LoadPublicKeys(data []byte) (map[string]string, error) {
	var pubKeys map[string]string
	err := json.Unmarshal(data, &pubKeys)
	if err != nil {
		return nil, err
	}
	return pubKeys, nil
}