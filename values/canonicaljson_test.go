package values

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CanonicalizeJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{name: "empty", input: []byte("{}"), want: []byte("{}")},
		{name: "simple", input: []byte(`{"a":1,"b":2}`), want: []byte(`{"a":1,"b":2}`)},
		{name: "complex", input: []byte(`{"a":1,"b":2,"c":{"d":3,"e":4}}`), want: []byte(`{"a":1,"b":2,"c":{"d":3,"e":4}}`)},
		{name: "nested", input: []byte(`{"a":1,"b":2,"c":{"d":3,"e":4,"f":{"g":5,"h":6}}}`), want: []byte(`{"a":1,"b":2,"c":{"d":3,"e":4,"f":{"g":5,"h":6}}}`)},
		{name: "numbers", input: []byte(`{"a":1.1,"b":2.2,"c":3.3}`), want: []byte(`{"a":1.1,"b":2.2,"c":3.3}`)},
		{name: "strings", input: []byte(`{"a":"hello","b":"world","c":"foo"}`), want: []byte(`{"a":"hello","b":"world","c":"foo"}`)},
		{name: "booleans", input: []byte(`{"a":true,"b":false,"c":true}`), want: []byte(`{"a":true,"b":false,"c":true}`)},
		{name: "nulls", input: []byte(`{"a":null,"b":null,"c":null}`), want: []byte(`{"a":null,"b":null,"c":null}`)},
		{name: "mixed", input: []byte(`{"a":1,"b":2,"c":{"d":3,"e":4,"f":{"g":5,"h":6}}}`), want: []byte(`{"a":1,"b":2,"c":{"d":3,"e":4,"f":{"g":5,"h":6}}}`)},
		{name: "maps", input: []byte(`{"b":[2,1,3],"a":["2","1","3"],"c": {"d":[3,4],"e":4,"f":{"g":5,"h":6 } } }`), want: []byte(`{"a":["2","1","3"],"b":[2,1,3],"c":{"d":[3,4],"e":4,"f":{"g":5,"h":6}}}`)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := CanonicalizeJSON(test.input)
			require.NoError(t, err)
			assert.Equal(t, string(test.want), string(got))
		})
	}
}

func Test_Canonicalize_Mapany(t *testing.T) {
	m := MapAny{
		"d": MapAny{
			"f": MapAny{
				"h": 6,
				"g": 5,
			},
			"e": 4,
			"d": []int{3, 4},
		},
		"c": nil,
		"b": []int{2, 1, 3},
		"a": []string{"2", "1", "3"},
	}
	got, err := m.CanonicalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"a":["2","1","3"],"b":[2,1,3],"c":null,"d":{"d":[3,4],"e":4,"f":{"g":5,"h":6}}}`, string(got))

	idx := MapAny{
		"settings": MapAny{
			"index": json.RawMessage(`{
"number_of_shards":   1,
"number_of_replicas": 0
}`),
			"analysis": json.RawMessage(`{
"analyzer": {
	"security_text": {
		"type":      "custom",
		"tokenizer": "pattern",
		"pattern": "[\\s,;()\\[\\]{}<>\"']+",
		"filter": ["lowercase"]
	}
}
}`),
		},
		"mappings": MapAny{
			"dynamic_templates": json.RawMessage(`[
{
	"strings_as_keyword": {
		"match_mapping_type": "string",
		"mapping": {
			"type": "keyword",
			"ignore_above": 256
		}
	}
}
]`),
			"properties": MapAny{
				"type": MapAny{"type": "integer"},
				"BatchID": MapAny{
					"type":       "keyword",
					"doc_values": true,
				},
				"CloudID": MapAny{
					"type":       "keyword",
					"doc_values": true,
				},
				"batch_id": MapAny{
					"type": "alias",
					"path": "BatchID",
				},
				"cloud_id": MapAny{
					"type": "alias",
					"path": "CloudID",
				},
			},
		},
	}

	got, err = idx.CanonicalJSON()
	require.NoError(t, err)
	exp := `{"mappings":{"dynamic_templates":[{"strings_as_keyword":{"mapping":{"ignore_above":256,"type":"keyword"},"match_mapping_type":"string"}}],"properties":{"BatchID":{"doc_values":true,"type":"keyword"},"CloudID":{"doc_values":true,"type":"keyword"},"batch_id":{"path":"BatchID","type":"alias"},"cloud_id":{"path":"CloudID","type":"alias"},"type":{"type":"integer"}}},"settings":{"analysis":{"analyzer":{"security_text":{"filter":["lowercase"],"pattern":"[\\s,;()\\[\\]{}\u003c\u003e\"']+","tokenizer":"pattern","type":"custom"}}},"index":{"number_of_replicas":0,"number_of_shards":1}}}`
	assert.Equal(t, exp, string(got))

	var buf bytes.Buffer
	require.NoError(t, json.Indent(&buf, got, "", "  "))
	exp2 := `{
  "mappings": {
    "dynamic_templates": [
      {
        "strings_as_keyword": {
          "mapping": {
            "ignore_above": 256,
            "type": "keyword"
          },
          "match_mapping_type": "string"
        }
      }
    ],
    "properties": {
      "BatchID": {
        "doc_values": true,
        "type": "keyword"
      },
      "CloudID": {
        "doc_values": true,
        "type": "keyword"
      },
      "batch_id": {
        "path": "BatchID",
        "type": "alias"
      },
      "cloud_id": {
        "path": "CloudID",
        "type": "alias"
      },
      "type": {
        "type": "integer"
      }
    }
  },
  "settings": {
    "analysis": {
      "analyzer": {
        "security_text": {
          "filter": [
            "lowercase"
          ],
          "pattern": "[\\s,;()\\[\\]{}\u003c\u003e\"']+",
          "tokenizer": "pattern",
          "type": "custom"
        }
      }
    },
    "index": {
      "number_of_replicas": 0,
      "number_of_shards": 1
    }
  }
}`
	assert.Equal(t, exp2, buf.String())
}

func Test_MarshalCanonicalJSON(t *testing.T) {
	m := MapAny{
		"b": []int{2, 1, 3},
		"a": []string{"2", "1", "3"},
		"d": MapAny{
			"f": MapAny{
				"g":      5,
				"h":      6,
				"false":  false,
				"true":   true,
				"null":   nil,
				"string": "",
				"number": 1.1,
				"int":    0,
				"uint":   uint(0),
				"uint64": uint64(0),
				"int64":  int64(0),
			},
			"d": []int{3, 4},
			"e": 4,
		},
		"c":      nil,
		"e":      1.1,
		"false":  false,
		"true":   true,
		"null":   nil,
		"string": "string",
		"number": 1.1,
		"int":    1,
		"uint":   uint(1),
		"uint64": uint64(1),
		"int64":  int64(1),
	}
	got, err := MarshalCanonicalJSON(m)
	require.NoError(t, err)
	exp := `{"a":["2","1","3"],"b":[2,1,3],"c":null,"d":{"d":[3,4],"e":4,"f":{"false":false,"g":5,"h":6,"int":0,"int64":0,"null":null,"number":1.1,"string":"","true":true,"uint":0,"uint64":0}},"e":1.1,"false":false,"int":1,"int64":1,"null":null,"number":1.1,"string":"string","true":true,"uint":1,"uint64":1}`
	assert.Equal(t, exp, string(got))
}

// canonicalInner declares the fields in non-alphabetical order,
// as encoding/json preserves the declaration order.
type canonicalInner struct {
	Zulu  string `json:"zulu"`
	Alpha string `json:"alpha"`
}

type canonicalStruct struct {
	Zebra    string          `json:"zebra"`
	Alpha    int             `json:"alpha"`
	Inner    canonicalInner  `json:"inner"`
	Raw      json.RawMessage `json:"raw"`
	Ptr      *canonicalInner `json:"ptr"`
	Ignored  string          `json:"-"`
	Optional string          `json:"optional,omitempty"`
}

// unsortedMarshaler emits an object with unsorted keys.
type unsortedMarshaler struct{}

func (unsortedMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(`{"z":1,"a":{"y":2,"b":3}}`), nil
}

type failingMarshaler struct{}

func (failingMarshaler) MarshalJSON() ([]byte, error) {
	return nil, errors.New("failed to marshal")
}

func Test_MarshalCanonicalJSON_Types(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`{"z":1,"a":2}`)

	tests := []struct {
		name string
		val  any
		want string
	}{
		{name: "nil", val: nil, want: `null`},
		{name: "raw_object", val: json.RawMessage(` {"z" : 1, "a" : {"y":2, "b":3} } `), want: `{"a":{"b":3,"y":2},"z":1}`},
		{name: "raw_array", val: json.RawMessage(`[{"z":1,"a":2},[{"d":1,"c":2}]]`), want: `[{"a":2,"z":1},[{"c":2,"d":1}]]`},
		{name: "raw_scalar", val: json.RawMessage(`"str"`), want: `"str"`},
		{name: "raw_empty", val: json.RawMessage(``), want: `null`},
		{name: "raw_nil", val: json.RawMessage(nil), want: `null`},
		{name: "raw_null", val: json.RawMessage(`null`), want: `null`},
		{name: "raw_ptr", val: &raw, want: `{"a":2,"z":1}`},
		{name: "raw_ptr_nil", val: (*json.RawMessage)(nil), want: `null`},
		{name: "raw_big_number", val: json.RawMessage(`{"n":123456789012345678901234567890}`), want: `{"n":123456789012345678901234567890}`},
		{name: "raw_in_slice", val: []any{json.RawMessage(`{"z":1,"a":2}`)}, want: `[{"a":2,"z":1}]`},
		{name: "raw_slice", val: []json.RawMessage{json.RawMessage(`{"z":1,"a":2}`)}, want: `[{"a":2,"z":1}]`},
		{
			name: "struct",
			val: canonicalStruct{
				Zebra: "z",
				Alpha: 1,
				Inner: canonicalInner{Zulu: "z", Alpha: "a"},
				Raw:   json.RawMessage(`{"z":1,"a":2}`),
				Ptr:   &canonicalInner{Zulu: "z", Alpha: "a"},
			},
			want: `{"alpha":1,"inner":{"alpha":"a","zulu":"z"},"ptr":{"alpha":"a","zulu":"z"},"raw":{"a":2,"z":1},"zebra":"z"}`,
		},
		{
			name: "struct_ptr",
			val:  &canonicalStruct{Zebra: "z", Raw: json.RawMessage(`null`)},
			want: `{"alpha":0,"inner":{"alpha":"","zulu":""},"ptr":null,"raw":null,"zebra":"z"}`,
		},
		{name: "struct_in_map", val: MapAny{"s": canonicalInner{Zulu: "z", Alpha: "a"}}, want: `{"s":{"alpha":"a","zulu":"z"}}`},
		{name: "struct_slice", val: []canonicalInner{{Zulu: "z", Alpha: "a"}}, want: `[{"alpha":"a","zulu":"z"}]`},
		{name: "mapany_slice", val: []MapAny{{"z": 1, "a": MapAny{"y": 2, "b": 3}}}, want: `[{"a":{"b":3,"y":2},"z":1}]`},
		{name: "map_slice", val: []map[string]any{{"z": 1, "a": 2}}, want: `[{"a":2,"z":1}]`},
		{name: "map_string", val: map[string]string{"z": "1", "a": "2"}, want: `{"a":"2","z":"1"}`},
		{name: "map_mapany", val: map[string]MapAny{"z": {"y": 1, "b": 2}}, want: `{"z":{"b":2,"y":1}}`},
		{name: "map_of_raw", val: map[string]json.RawMessage{"z": json.RawMessage(`{"y":1,"b":2}`)}, want: `{"z":{"b":2,"y":1}}`},
		{name: "marshaler", val: unsortedMarshaler{}, want: `{"a":{"b":3,"y":2},"z":1}`},
		{name: "marshaler_in_map", val: MapAny{"m": unsortedMarshaler{}}, want: `{"m":{"a":{"b":3,"y":2},"z":1}}`},
		{name: "time", val: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC), want: `"2024-01-02T03:04:05Z"`},
		{name: "bytes", val: []byte("abc"), want: `"YWJj"`},
		{name: "html_escaped", val: "<a>", want: `"\u003ca\u003e"`},
		{name: "number", val: json.Number("1.100"), want: `1.100`},
		{name: "number_empty", val: json.Number(""), want: `0`},
		{name: "empty_map", val: MapAny{}, want: `{}`},
		{name: "empty_slice", val: []any{}, want: `[]`},
		{name: "nil_mapany", val: MapAny(nil), want: `{}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := MarshalCanonicalJSON(test.val)
			require.NoError(t, err)
			assert.Equal(t, test.want, string(got))
			assert.True(t, json.Valid(got), "invalid JSON: %s", got)

			// the canonical form must be stable
			again, err := CanonicalizeJSON(got)
			require.NoError(t, err)
			assert.Equal(t, test.want, string(again))
		})
	}
}

func Test_MarshalCanonicalJSON_Errors(t *testing.T) {
	t.Parallel()

	deep := MapAny{}
	cur := deep
	for range maxCanonicalJSONDepth + 1 {
		next := MapAny{}
		cur["n"] = next
		cur = next
	}

	cyclic := MapAny{}
	cyclic["self"] = cyclic

	tests := []struct {
		name   string
		val    any
		experr string
	}{
		{name: "invalid_raw", val: json.RawMessage(`{"a":`), experr: "decode json: unexpected EOF"},
		{name: "invalid_raw_token", val: json.RawMessage(`{a:1}`), experr: "decode json: invalid character 'a' looking for beginning of object key string"},
		{name: "trailing_data", val: json.RawMessage(`{"a":1} {"b":2}`), experr: "decode json: unexpected data after top-level value"},
		{name: "trailing_garbage", val: json.RawMessage(`{"a":1} }`), experr: "decode json: invalid character '}' looking for beginning of value"},
		{name: "invalid_raw_in_map", val: MapAny{"a": json.RawMessage(`}`)}, experr: "decode json: invalid character '}' looking for beginning of value"},
		{name: "invalid_raw_in_slice", val: []any{1, json.RawMessage(`}`)}, experr: "decode json: invalid character '}' looking for beginning of value"},
		{name: "unsupported_value", val: MapAny{"a": make(chan int)}, experr: "marshal value: json: unsupported type: chan int"},
		{name: "failing_marshaler", val: failingMarshaler{}, experr: "marshal value: json: error calling MarshalJSON for type *values.failingMarshaler: failed to marshal"},
		{name: "too_deep", val: deep, experr: "json value exceeds max nesting depth: 1000"},
		{name: "cyclic", val: cyclic, experr: "json value exceeds max nesting depth: 1000"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := MarshalCanonicalJSON(test.val)
			require.Error(t, err)
			assert.Equal(t, test.experr, err.Error())
		})
	}
}

func Test_CanonicalizeJSON_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  []byte
		want   string
		experr string
	}{
		{name: "empty", input: nil, experr: "decode json: empty input"},
		{name: "blank", input: []byte("  \n"), experr: "decode json: empty input"},
		{name: "scalar", input: []byte(`  "a"  `), want: `"a"`},
		{name: "invalid", input: []byte(`{`), experr: "decode json: unexpected EOF"},
		{name: "trailing", input: []byte(`{} {}`), experr: "decode json: unexpected data after top-level value"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := CanonicalizeJSON(test.input)
			if test.experr != "" {
				require.Error(t, err)
				assert.Equal(t, test.experr, err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, string(got))
		})
	}
}
