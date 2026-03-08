package main

import (
	"fmt"
	"os"

	"github.com/dip-protocol/dip-go-verifier/verifier"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: verify artifact.json")
		return
	}

	artifactFile := os.Args[1]

	ok, err := verifier.VerifyArtifact(artifactFile)

	if err != nil {
		fmt.Println("Verification error:", err)
		return
	}

	if !ok {
		fmt.Println("Verification failed")
		return
	}

	fmt.Println("Artifact verified successfully")
}