package signature

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
)

// VerifySignature verifies that sigHex is a valid signature of data
// using the public key pubKeyHex (both hex-encoded).
func VerifySignature(pubKeyHex, sigHex string, data []byte) (bool, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return false, err
	}

	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		return false, err
	}

	if len(pubKeyBytes) != ed25519.PublicKeySize {
		return false, errors.New("invalid public key size")
	}
	if len(sigBytes) != ed25519.SignatureSize {
		return false, errors.New("invalid signature size")
	}

	valid := ed25519.Verify(pubKeyBytes, data, sigBytes)
	return valid, nil
}