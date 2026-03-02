package artifact

import "encoding/json"

type DecisionArtifact struct {
    ArtifactID   string `json:"artifact_id"`
    ArtifactType string `json:"artifact_type"`
    Signature    string `json:"signature"`
}

func ParseDecisionArtifact(data []byte) (*DecisionArtifact, error) {
    var da DecisionArtifact
    err := json.Unmarshal(data, &da)
    if err != nil {
        return nil, err
    }
    return &da, nil
}