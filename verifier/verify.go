package verifier

import (
	"crypto/ed25519"
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

	canonical, err := CanonicalizeJSON(decisionBytes)
	if err != nil {
		return false, err
	}

	expectedID := ComputeArtifactID(canonical)

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