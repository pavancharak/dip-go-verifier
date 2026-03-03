package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pavancharak/dip-go-verifier/internal/artifact"
	"github.com/pavancharak/dip-go-verifier/internal/canonicalization"
	"github.com/pavancharak/dip-go-verifier/internal/hashing"
	"github.com/pavancharak/dip-go-verifier/internal/policy"
	"github.com/pavancharak/dip-go-verifier/internal/signature"
)

func main() {

	artifactPath := flag.String("artifact", "", "Path to artifact JSON file")
	pubKeysPath := flag.String("pubkeys", "", "Path to public keys JSON file")

	flag.Parse()

	if *artifactPath == "" || *pubKeysPath == "" {
		fmt.Println("Usage:")
		fmt.Println("  dip-verify --artifact <artifact.json> --pubkeys <public_keys.json>")
		os.Exit(1)
	}

	// Step 1: Read artifact
	data, err := ioutil.ReadFile(*artifactPath)
	if err != nil {
		log.Fatalf("failed to read artifact file: %v", err)
	}

	// Step 2: Parse artifact
	da, err := artifact.ParseDecisionArtifact(data)
	if err != nil {
		log.Fatalf("artifact parsing failed: %v", err)
	}

	// Step 3: Strict policy validation
	if err := policy.ValidateProtocolVersion(da.ProtocolVersion); err != nil {
		log.Fatalf("protocol validation failed: %v", err)
	}

	if err := policy.ValidateSignatureAlgorithm(da.SignatureAlgorithm); err != nil {
		log.Fatalf("algorithm validation failed: %v", err)
	}

	if err := policy.ValidateTimestamp(da.TimestampUTC); err != nil {
		log.Fatalf("timestamp validation failed: %v", err)
	}

	if err := policy.ValidateLowerHex(da.Signature); err != nil {
		log.Fatalf("signature hex validation failed: %v", err)
	}

	if err := policy.ValidateLowerHex(da.DecisionHash); err != nil {
		log.Fatalf("decision_hash validation failed: %v", err)
	}

	if err := policy.ValidateLowerHex(da.InputHash); err != nil {
		log.Fatalf("input_hash validation failed: %v", err)
	}

	if err := policy.ValidateLowerHex(da.AuthorityHash); err != nil {
		log.Fatalf("authority_hash validation failed: %v", err)
	}

	// Step 4: Canonicalize (signature excluded by canonical layer)
	canonical := canonicalization.Canonicalize(data)

	fmt.Println("Canonical JSON:")
	fmt.Println(string(canonical))
	fmt.Println()

	fmt.Println("Canonical Bytes (hex):")
	fmt.Println(hex.EncodeToString(canonical))
	fmt.Println()

	hash := hashing.ComputeSHA256(canonical)
	fmt.Println("Canonical SHA-256:", hash)
	fmt.Println()

	// Step 5: Load public keys
	pubKeysBytes, err := ioutil.ReadFile(*pubKeysPath)
	if err != nil {
		log.Fatalf("failed to read public keys file: %v", err)
	}

	var pubKeys map[string]string
	err = json.Unmarshal(pubKeysBytes, &pubKeys)
	if err != nil {
		log.Fatalf("invalid public keys JSON: %v", err)
	}

	pubKeyHex, ok := pubKeys[da.ArtifactID]
	if !ok {
		log.Fatalf("public key not found for artifact %s", da.ArtifactID)
	}

	if err := policy.ValidateLowerHex(pubKeyHex); err != nil {
		log.Fatalf("public key hex validation failed: %v", err)
	}

	// Step 6: Verify signature
	valid, err := signature.VerifySignature(pubKeyHex, da.Signature, canonical)
	if err != nil {
		log.Fatalf("signature verification error: %v", err)
	}

	if !valid {
		fmt.Println("[FAIL] Signature verification failed")
		os.Exit(1)
	}

	fmt.Println("[PASS] Signature verified")
}