package values

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/cockroachdb/errors"
)

func MarshalCanonicalJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := WriteCanonicalJSON(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Canonicalize returns the canonical JSON representation of the input
func CanonicalizeJSON(input []byte) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(input))
	dec.UseNumber()

	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, errors.Wrap(err, "decode json")
	}

	var buf bytes.Buffer
	if err := WriteCanonicalJSON(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func WriteCanonicalJSON(buf *bytes.Buffer, v any) error {
	switch x := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		buf.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			// key
			kb, _ := json.Marshal(k)
			buf.Write(kb)
			buf.WriteByte(':')

			// value
			if err := WriteCanonicalJSON(buf, x[k]); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
		return nil

	case []any:
		buf.WriteByte('[')
		for i := range x {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := WriteCanonicalJSON(buf, x[i]); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
		return nil

	case json.Number:
		// Preserve exact lexical form as provided by decoder
		buf.WriteString(x.String())
		return nil

	case string, bool, nil, float64:
		b, err := json.Marshal(x)
		if err != nil {
			return errors.Wrap(err, "marshal scalar")
		}
		buf.Write(b)
		return nil

	default:
		// should be rare, but handle safely
		b, err := json.Marshal(x)
		if err != nil {
			return errors.Wrap(err, "marshal fallback")
		}
		buf.Write(b)
		return nil
	}
}
