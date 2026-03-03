package canonicalization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// Canonicalize returns canonical JSON bytes excluding the "signature" field.
func Canonicalize(data []byte) []byte {

	var obj map[string]interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return nil
	}

	// Remove signature field from preimage
	delete(obj, "signature")

	var buf bytes.Buffer
	writeCanonical(&buf, obj)

	return buf.Bytes()
}

func writeCanonical(buf *bytes.Buffer, v interface{}) {

	switch val := v.(type) {

	case map[string]interface{}:
		buf.WriteByte('{')

		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return bytes.Compare([]byte(keys[i]), []byte(keys[j])) < 0
		})

		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			writeString(buf, k)
			buf.WriteByte(':')
			writeCanonical(buf, val[k])
		}

		buf.WriteByte('}')

	case []interface{}:
		buf.WriteByte('[')
		for i, elem := range val {
			if i > 0 {
				buf.WriteByte(',')
			}
			writeCanonical(buf, elem)
		}
		buf.WriteByte(']')

	case string:
		writeString(buf, val)

	case float64:
		// JSON numbers decode as float64
		// DIP prohibits floats — assume integer only
		if val != float64(int64(val)) {
			// float detected (not allowed by DIP)
			buf.WriteString("INVALID_FLOAT")
			return
		}
		buf.WriteString(fmt.Sprintf("%d", int64(val)))

	case bool:
		if val {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}

	case nil:
		// Null not allowed in DIP canonical form
		buf.WriteString("null")

	default:
		// Fallback safe marshal
		b, _ := json.Marshal(val)
		buf.Write(b)
	}
}

func writeString(buf *bytes.Buffer, s string) {
	b, _ := json.Marshal(s)
	buf.Write(b)
}