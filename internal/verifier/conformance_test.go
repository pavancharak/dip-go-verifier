package verifier

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pavancharak/dip-go-verifier/internal/artifact"
	"github.com/pavancharak/dip-go-verifier/internal/canonicalization"
	"github.com/pavancharak/dip-go-verifier/internal/hashing"
	"github.com/pavancharak/dip-go-verifier/internal/signature"
)

func loadPublicKeys(t *testing.T) map[string]string {
	pubKeysBytes, err := ioutil.ReadFile("../../testdata/vectors/public_keys.json")
	if err != nil {
		t.Fatal(err)
	}

	var pubKeys map[string]string
	err = json.Unmarshal(pubKeysBytes, &pubKeys)
	if err != nil {
		t.Fatal(err)
	}

	return pubKeys
}

func TestConformanceVectors(t *testing.T) {

	pubKeys := loadPublicKeys(t)

	err := filepath.Walk("../../testdata/vectors", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		if filepath.Base(path) == "public_keys.json" {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		da, err := artifact.ParseDecisionArtifact(data)

		isValidVector := filepath.Base(path) == "valid_decision.json"

		// Parsing failure is expected for invalid vectors
		if err != nil {
			if isValidVector {
				t.Fatalf("Valid vector failed parsing: %s", path)
			}
			return nil
		}

		canonical := canonicalization.Canonicalize(data)
		hashHex := hashing.ComputeSHA256(canonical)

		hashBytes, err := hex.DecodeString(hashHex)
		if err != nil {
			t.Fatal(err)
		}

		pubKeyHex, ok := pubKeys[da.ArtifactID]
		if !ok {
			if isValidVector {
				t.Fatalf("Public key missing for valid vector: %s", path)
			}
			return nil
		}

		valid, err := signature.VerifySignature(pubKeyHex, da.Signature, hashBytes)

		// Verification error is expected for invalid vectors
		if err != nil {
			if isValidVector {
				t.Fatalf("Valid vector failed verification: %s", path)
			}
			return nil
		}

		if isValidVector && !valid {
			t.Fatalf("Expected valid vector to pass: %s", path)
		}

		if !isValidVector && valid {
			t.Fatalf("Invalid vector incorrectly passed: %s", path)
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}