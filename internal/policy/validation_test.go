package policy

import "testing"

func TestValidateLowerHex(t *testing.T) {
	valid := "aabbccddeeff00112233445566778899"
	if err := ValidateLowerHex(valid); err != nil {
		t.Fatalf("expected valid hex, got error: %v", err)
	}

	upper := "AABBCC"
	if err := ValidateLowerHex(upper); err == nil {
		t.Fatalf("expected uppercase hex to fail")
	}

	oddLength := "abc"
	if err := ValidateLowerHex(oddLength); err == nil {
		t.Fatalf("expected odd-length hex to fail")
	}

	invalidChar := "zz"
	if err := ValidateLowerHex(invalidChar); err == nil {
		t.Fatalf("expected invalid hex characters to fail")
	}
}

func TestValidateTimestamp(t *testing.T) {
	valid := "2026-01-01T00:00:00Z"
	if err := ValidateTimestamp(valid); err != nil {
		t.Fatalf("expected valid timestamp, got error: %v", err)
	}

	noZ := "2026-01-01T00:00:00"
	if err := ValidateTimestamp(noZ); err == nil {
		t.Fatalf("expected timestamp without Z to fail")
	}

	invalid := "not-a-time"
	if err := ValidateTimestamp(invalid); err == nil {
		t.Fatalf("expected invalid timestamp to fail")
	}
}

func TestValidateSignatureAlgorithm(t *testing.T) {
	if err := ValidateSignatureAlgorithm("Ed25519"); err != nil {
		t.Fatalf("expected Ed25519 to pass")
	}

	if err := ValidateSignatureAlgorithm("RSA"); err == nil {
		t.Fatalf("expected non-Ed25519 algorithm to fail")
	}
}

func TestValidateProtocolVersion(t *testing.T) {
	if err := ValidateProtocolVersion("1"); err != nil {
		t.Fatalf("expected protocol version 1 to pass")
	}

	if err := ValidateProtocolVersion("2"); err == nil {
		t.Fatalf("expected protocol version != 1 to fail")
	}
}