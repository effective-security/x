package values

import (
	"testing"

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
		"b": []int{2, 1, 3},
		"a": []string{"2", "1", "3"},
		"d": MapAny{
			"d": []int{3, 4},
			"e": 4,
			"f": MapAny{
				"g": 5,
				"h": 6,
			},
		},
		"c": nil,
	}
	got, err := m.CanonicalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"a":["2","1","3"],"b":[2,1,3],"c":null,"d":{"d":[3,4],"e":4,"f":{"g":5,"h":6}}}`, string(got))
}
