package canonicalization

import (
	"encoding/hex"
	"io/ioutil"
	"testing"

	"github.com/pavancharak/dip-go-verifier/internal/hashing"
)

func TestCanonicalDeterminism(t *testing.T) {

	data, err := ioutil.ReadFile("../../testdata/vectors/valid_decision.json")
	if err != nil {
		t.Fatal(err)
	}

	canonical := Canonicalize(data)

	hashHex := hashing.ComputeSHA256(canonical)

	expected := "b84a240b0a5e48ba8248799ae9c9bf91f8482ed7ca2a16ec0403a00ea315ed28"

	if hashHex != expected {
		t.Fatalf("Canonical SHA-256 mismatch.\nExpected: %s\nGot:      %s\nCanonical (hex): %s",
			expected,
			hashHex,
			hex.EncodeToString(canonical),
		)
	}
}