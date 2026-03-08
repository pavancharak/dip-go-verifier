package verifier

import "testing"

func TestVerifyArtifact(t *testing.T) {

	ok, err := VerifyArtifact("../testdata/artifact.json")

	if err != nil {
		t.Fatalf("verification returned error: %v", err)
	}

	if !ok {
		t.Fatalf("artifact verification failed")
	}
}