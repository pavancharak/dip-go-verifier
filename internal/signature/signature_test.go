package signature

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/pavancharak/dip-go-verifier/internal/artifact"
	"github.com/pavancharak/dip-go-verifier/internal/canonicalization"
	"github.com/pavancharak/dip-go-verifier/internal/hashing"
)

func loadArtifact(t *testing.T) ([]byte, *artifact.DecisionArtifact, string) {

	data, err := ioutil.ReadFile("../../testdata/vectors/valid_decision.json")
	if err != nil {
		t.Fatal(err)
	}

	da, err := artifact.ParseDecisionArtifact(data)
	if err != nil {
		t.Fatal(err)
	}

	pubKeysBytes, err := ioutil.ReadFile("../../testdata/vectors/public_keys.json")
	if err != nil {
		t.Fatal(err)
	}

	var pubKeys map[string]string
	err = json.Unmarshal(pubKeysBytes, &pubKeys)
	if err != nil {
		t.Fatal(err)
	}

	pubKeyHex := pubKeys[da.ArtifactID]

	return data, da, pubKeyHex
}

func TestSignatureVerification(t *testing.T) {

	data, da, pubKeyHex := loadArtifact(t)

	canonical := canonicalization.Canonicalize(data)
	hashHex := hashing.ComputeSHA256(canonical)
	hashBytes, err := hex.DecodeString(hashHex)
	if err != nil {
		t.Fatal(err)
	}

	// --- Positive test (must pass) ---
	valid, err := VerifySignature(pubKeyHex, da.Signature, hashBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("Expected valid signature but verification failed")
	}

	// --- Negative test 1: tampered signature ---
	tamperedSig := da.Signature[:len(da.Signature)-2] + "00"

	valid, err = VerifySignature(pubKeyHex, tamperedSig, hashBytes)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("Expected signature verification to fail for tampered signature")
	}

	// --- Negative test 2: tampered canonical JSON ---
	canonical[0] ^= 0xFF
	hashHex = hashing.ComputeSHA256(canonical)
	hashBytes, err = hex.DecodeString(hashHex)
	if err != nil {
		t.Fatal(err)
	}

	valid, err = VerifySignature(pubKeyHex, da.Signature, hashBytes)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("Expected signature verification to fail for tampered data")
	}
}