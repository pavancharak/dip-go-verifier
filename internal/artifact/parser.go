package artifact

import (
	"encoding/json"
)

type DecisionArtifact struct {
	ArtifactID          string `json:"artifact_id"`
	ArtifactType        string `json:"artifact_type"`
	ArtifactVersion     int    `json:"artifact_version"`
	ProtocolVersion     string `json:"protocol_version"`
	SignatureAlgorithm  string `json:"signature_algorithm"`
	TimestampUTC        string `json:"timestamp_utc"`

	DecisionHash        string `json:"decision_hash"`
	InputHash           string `json:"input_hash"`
	AuthorityHash       string `json:"authority_hash"`
	AuthorityLevel      string `json:"authority_level"`

	ModelID             string `json:"model_id"`
	ModelVersion        string `json:"model_version"`
	ProducerID          string `json:"producer_id"`
	ProducerPublicKeyID string `json:"producer_public_key_id"`

	Signature           string `json:"signature"`
}

func ParseDecisionArtifact(data []byte) (*DecisionArtifact, error) {
	var da DecisionArtifact
	err := json.Unmarshal(data, &da)
	if err != nil {
		return nil, err
	}
	return &da, nil
}