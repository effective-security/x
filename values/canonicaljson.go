package values

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/cockroachdb/errors"
)

// maxCanonicalJSONDepth limits the nesting level of the encoded value,
// to protect against cyclic values that otherwise exhaust the stack.
const maxCanonicalJSONDepth = 1000

// MarshalCanonicalJSON returns the canonical JSON representation of the value,
// with all object keys sorted alphabetically.
func MarshalCanonicalJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := WriteCanonicalJSON(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// CanonicalizeJSON returns the canonical JSON representation of the input,
// with all object keys sorted alphabetically.
func CanonicalizeJSON(input []byte) ([]byte, error) {
	if len(bytes.TrimSpace(input)) == 0 {
		return nil, errors.New("decode json: empty input")
	}

	var buf bytes.Buffer
	if err := writeCanonicalRawJSON(&buf, input, 0); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WriteCanonicalJSON writes the canonical JSON representation of the input to the buffer.
// All keys are sorted alphabetically, at any nesting level.
//
// Values that are not handled natively, such as structs, typed maps and slices,
// json.RawMessage or custom json.Marshaler implementations, are encoded with
// encoding/json and then re-encoded from the resulting JSON,
// as encoding/json preserves the declaration order of struct fields
// and the original order of raw JSON.
func WriteCanonicalJSON(buf *bytes.Buffer, v any) error {
	return writeCanonicalJSON(buf, v, 0)
}

func writeCanonicalJSON(buf *bytes.Buffer, v any, depth int) error {
	if depth > maxCanonicalJSONDepth {
		return errors.Errorf("json value exceeds max nesting depth: %d", maxCanonicalJSONDepth)
	}

	switch x := v.(type) {
	case nil:
		buf.WriteString("null")
		return nil

	case MapAny:
		return writeCanonicalObject(buf, x, depth)

	case map[string]any:
		return writeCanonicalObject(buf, MapAny(x), depth)

	case []any:
		buf.WriteByte('[')
		for i := range x {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := writeCanonicalJSON(buf, x[i], depth+1); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
		return nil

	case json.RawMessage:
		return writeCanonicalRawJSON(buf, x, depth)

	case *json.RawMessage:
		if x == nil {
			buf.WriteString("null")
			return nil
		}
		return writeCanonicalRawJSON(buf, *x, depth)

	case string, bool,
		float32, float64,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		// json.Number preserves the exact lexical form as provided by the decoder
		json.Number:
		b, err := json.Marshal(x)
		if err != nil {
			return errors.Wrap(err, "marshal scalar")
		}
		buf.Write(b)
		return nil

	default:
		b, err := json.Marshal(x)
		if err != nil {
			return errors.Wrap(err, "marshal value")
		}
		return writeCanonicalRawJSON(buf, b, depth)
	}
}

func writeCanonicalObject(buf *bytes.Buffer, m MapAny, depth int) error {
	keys := m.OrderedKeys()

	buf.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		kb, err := json.Marshal(k)
		if err != nil {
			return errors.Wrap(err, "marshal key")
		}
		buf.Write(kb)
		buf.WriteByte(':')

		if err := writeCanonicalJSON(buf, m[k], depth+1); err != nil {
			return err
		}
	}
	buf.WriteByte('}')
	return nil
}

// writeCanonicalRawJSON decodes already encoded JSON and writes it back in canonical form.
func writeCanonicalRawJSON(buf *bytes.Buffer, raw []byte, depth int) error {
	if len(bytes.TrimSpace(raw)) == 0 {
		buf.WriteString("null")
		return nil
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()

	var v any
	if err := dec.Decode(&v); err != nil {
		return errors.Wrap(err, "decode json")
	}
	if err := dec.Decode(new(json.RawMessage)); !errors.Is(err, io.EOF) {
		if err != nil {
			return errors.Wrap(err, "decode json")
		}
		return errors.New("decode json: unexpected data after top-level value")
	}

	return writeCanonicalJSON(buf, v, depth)
}
