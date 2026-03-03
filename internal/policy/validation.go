package policy

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidHex       = errors.New("invalid lowercase hex string")
	ErrInvalidTimestamp = errors.New("invalid RFC3339 UTC timestamp")
	ErrInvalidAlgorithm = errors.New("unsupported signature algorithm")
	ErrInvalidProtocol  = errors.New("unsupported protocol version")
)

// ValidateLowerHex ensures string is lowercase hex and even length
func ValidateLowerHex(s string) error {
	if len(s)%2 != 0 {
		return ErrInvalidHex
	}
	if s != strings.ToLower(s) {
		return ErrInvalidHex
	}
	_, err := hex.DecodeString(s)
	if err != nil {
		return ErrInvalidHex
	}
	return nil
}

// ValidateTimestamp ensures strict RFC3339 UTC ending in Z
func ValidateTimestamp(ts string) error {
	if !strings.HasSuffix(ts, "Z") {
		return ErrInvalidTimestamp
	}
	_, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ErrInvalidTimestamp
	}
	return nil
}

// ValidateSignatureAlgorithm enforces Ed25519 in v1
func ValidateSignatureAlgorithm(alg string) error {
	if alg != "Ed25519" {
		return ErrInvalidAlgorithm
	}
	return nil
}

// ValidateProtocolVersion enforces v1 only
func ValidateProtocolVersion(version string) error {
	if version != "1" {
		return ErrInvalidProtocol
	}
	return nil
}