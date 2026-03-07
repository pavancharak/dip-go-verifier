package verifier

import (
	"bytes"
	"encoding/json"
)

func CanonicalizeJSON(data []byte) ([]byte, error) {

	var obj any

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err = enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(buf.Bytes()), nil
}