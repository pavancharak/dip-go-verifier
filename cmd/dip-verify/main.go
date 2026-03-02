package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pavancharak/dip-go-verifier/internal/artifact"
	"github.com/pavancharak/dip-go-verifier/internal/canonicalization"
	"github.com/pavancharak/dip-go-verifier/internal/hashing"
	"github.com/pavancharak/dip-go-verifier/internal/signature"
)

func main() {
	// Step 1: Read artifact JSON
	data, err := ioutil.ReadFile("testdata/vectors/valid_decision.json")
	if err != nil {
		log.Fatal(err)
	}

	// Step 2: Parse artifact
	da, err := artifact.ParseDecisionArtifact(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Artifact ID:", da.ArtifactID)
	fmt.Println("Artifact Type:", da.ArtifactType)
	fmt.Println("Signature:", da.Signature)

	// Step 3: Canonicalize and compute hash
	canonical := canonicalization.Canonicalize(data)
	hash := hashing.ComputeSHA256(canonical)
	fmt.Println("Canonical SHA-256 Hash:", hash)

	// Step 4: Load public keys
	pubKeysBytes, err := ioutil.ReadFile("testdata/vectors/public_keys.json")
	if err != nil {
		log.Fatal(err)
	}

	var pubKeys map[string]string
	err = json.Unmarshal(pubKeysBytes, &pubKeys)
	if err != nil {
		log.Fatal(err)
	}

	pubKeyHex, ok := pubKeys[da.ArtifactID]
	if !ok {
		log.Fatalf("public key not found for artifact %s", da.ArtifactID)
	}

	// Step 5: Verify signature
	valid, err := signature.VerifySignature(pubKeyHex, da.Signature, canonical)
	if err != nil {
		log.Fatal(err)
	}

	if valid {
		fmt.Println("[PASS] Signature verified")
	} else {
		fmt.Println("[FAIL] Signature verification failed")
	}
}