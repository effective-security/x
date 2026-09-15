package values

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// TestMapAnyScanUnsupportedType covers Scan default branch for unsupported types.
func TestMapAnyScanUnsupportedType(t *testing.T) {
	t.Parallel()
	var m MapAny
	err := m.Scan(123)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported scan type:")
}

// TestIndentJSONFallback covers IndentJSON error path when input is not valid JSON.
func TestIndentJSONFallback(t *testing.T) {
	t.Parallel()
	orig := "not json"
	got := IndentJSON(orig)
	assert.Equal(t, orig, got)
}

// TestMapAnyScanInvalidJSON covers Scan path when JSON is invalid and returns unmarshal error.
func TestMapAnyScanInvalidJSON(t *testing.T) {
	t.Parallel()
	var m MapAny
	// string value that is invalid JSON
	err := m.Scan("{invalid}")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestMapAnyFromProtoValue(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		FromStruct(nil)
	})
	assert.Panics(t, func() {
		FromStruct("string")
	})

	type styp struct {
		S    string         `json:"field"`
		I    int            `json:"i"`
		Ob   MapAny         `json:"ob"`
		M    map[string]any `json:"m"`
		List []string       `json:"list"`
		Ints []int          `json:"ints"`

		priv string `json:"-"`
	}

	tcases := []struct {
		name string
		val  any
		exp  string
	}{
		{
			name: "basic struct",
			val: styp{
				S: "value",
				M: map[string]any{
					"field":  nil,
					"field2": []string{"", ""},
					"field3": []int{0, 0},
				},
			},
			exp: `{
  "S": "value"
}`,
		},
		{
			name: "ptr with private struct",
			val: &styp{
				S: "value",
				I: 1,
				Ob: MapAny{
					"field": "value",
				},
				M: map[string]any{
					"field": "value",
				},
				List: []string{"value", "", "value2"},
				Ints: []int{1, 2, 3, 0, 0, 0},
			},
			exp: `{
  "I": 1,
  "Ints": [
    1,
    2,
    3
  ],
  "List": [
    "value",
    "value2"
  ],
  "M": {
    "field": "value"
  },
  "Ob": {
    "field": "value"
  },
  "S": "value"
}`,
		},
	}

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			many := FromStruct(tt.val)

			var val map[string]any
			many.To(&val)

			ps, err := structpb.NewStruct(val)
			require.NoError(t, err)

			sv := MapAny(ps.AsMap()).Shrink()
			assert.Equal(t, tt.exp, sv.JSONIndent())
		})
	}
}

func TestMapAny_AddKeepsFalseAndZero(t *testing.T) {
	t.Parallel()
	m := MapAny{}.Add("ok", false).Add("n", 0).Add("s", "").Add("skip", nil)
	assert.Equal(t, false, m["ok"])
	assert.Equal(t, 0, m["n"])
	_, hasEmpty := m["s"]
	assert.False(t, hasEmpty)
	_, hasNil := m["skip"]
	assert.False(t, hasNil)
}

func TestMapAny_SliceNoPanic(t *testing.T) {
	t.Parallel()
	m := MapAny{
		"any":  []any{1, 2},
		"strs": []string{"a", "b"},
		"bad":  "nope",
	}
	assert.Equal(t, []any{1, 2}, m.Slice("any"))
	assert.Equal(t, []any{"a", "b"}, m.Slice("strs"))
	assert.Nil(t, m.Slice("bad"))
	assert.Nil(t, m.Slice("missing"))
	assert.Nil(t, MapAny(nil).Slice("x"))
}
