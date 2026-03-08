package verifier

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
)

type Artifact struct {
	ArtifactVersion string    `json:"artifact_version"`
	ArtifactID      string    `json:"artifact_id"`
	Decision        any       `json:"decision"`
	Signature       Signature `json:"signature"`
}

type Signature struct {
	Algorithm string `json:"algorithm"`
	PublicKey []byte `json:"public_key"`
	Value     []byte `json:"value"`
}

func canonicalizeJSON(data []byte) ([]byte, error) {

	var obj any

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err = enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(buf.Bytes()), nil
}

func computeArtifactID(canonical []byte) string {

	hash := sha256.Sum256(canonical)

	return hex.EncodeToString(hash[:])
}

func VerifyArtifact(file string) (bool, error) {

	data, err := os.ReadFile(file)
	if err != nil {
		return false, err
	}

	var artifact Artifact

	err = json.Unmarshal(data, &artifact)
	if err != nil {
		return false, errors.New("invalid artifact format")
	}

	decisionBytes, err := json.Marshal(artifact.Decision)
	if err != nil {
		return false, err
	}

	canonical, err := canonicalizeJSON(decisionBytes)
	if err != nil {
		return false, err
	}

	expectedID := computeArtifactID(canonical)

	if expectedID != artifact.ArtifactID {
		return false, errors.New("artifact_id mismatch")
	}

	valid := ed25519.Verify(
		artifact.Signature.PublicKey,
		canonical,
		artifact.Signature.Value,
	)

	return valid, nil
}