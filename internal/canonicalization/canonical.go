package canonicalization

// For now, canonicalization is just identity (later we implement field ordering, whitespace normalization, etc.)
func Canonicalize(data []byte) []byte {
	return data
}