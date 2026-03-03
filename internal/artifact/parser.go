package artifact

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type DecisionArtifact struct {
	ArtifactID          string `json:"artifact_id"`
	ArtifactType        string `json:"artifact_type"`
	ArtifactVersion     int    `json:"artifact_version"`
	AuthorityHash       string `json:"authority_hash"`
	AuthorityLevel      string `json:"authority_level"`
	DecisionHash        string `json:"decision_hash"`
	InputHash           string `json:"input_hash"`
	ModelID             string `json:"model_id"`
	ModelVersion        string `json:"model_version"`
	ProducerID          string `json:"producer_id"`
	ProducerPublicKeyID string `json:"producer_public_key_id"`
	ProtocolVersion     string `json:"protocol_version"`
	SignatureAlgorithm  string `json:"signature_algorithm"`
	TimestampUTC        string `json:"timestamp_utc"`
	Signature           string `json:"signature"`
}

func validateHexField(value string, expectedLength int) error {
	if len(value) != expectedLength {
		return errors.New("invalid hex length")
	}
	if strings.ToLower(value) != value {
		return errors.New("hex must be lowercase")
	}
	_, err := hex.DecodeString(value)
	if err != nil {
		return errors.New("invalid hex encoding")
	}
	return nil
}

func ParseDecisionArtifact(data []byte) (*DecisionArtifact, error) {

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	var da DecisionArtifact
	if err := decoder.Decode(&da); err != nil {
		return nil, err
	}

	// Required field validation
	if da.ArtifactID == "" ||
		da.ArtifactType == "" ||
		da.ProtocolVersion == "" ||
		da.SignatureAlgorithm == "" ||
		da.Signature == "" {
		return nil, errors.New("missing required fields")
	}

	// Protocol enforcement
	if da.ProtocolVersion != "1" {
		return nil, errors.New("unsupported protocol version")
	}

	if da.SignatureAlgorithm != "Ed25519" {
		return nil, errors.New("unsupported signature algorithm")
	}

	if da.ArtifactType != "decision" {
		return nil, errors.New("invalid artifact type")
	}

	// Strict hex validation
	if err := validateHexField(da.AuthorityHash, 64); err != nil {
		return nil, err
	}
	if err := validateHexField(da.DecisionHash, 64); err != nil {
		return nil, err
	}
	if err := validateHexField(da.InputHash, 64); err != nil {
		return nil, err
	}
	if err := validateHexField(da.ProducerPublicKeyID, 64); err != nil {
		return nil, err
	}
	if err := validateHexField(da.Signature, 128); err != nil {
		return nil, err
	}

	// Timestamp validation (RFC3339 UTC)
	t, err := time.Parse(time.RFC3339, da.TimestampUTC)
	if err != nil {
		return nil, errors.New("invalid timestamp format")
	}
	if t.Location() != time.UTC {
		return nil, errors.New("timestamp must be UTC")
	}

	return &da, nil
}