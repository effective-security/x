package values

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValues(t *testing.T) {
	c := MapAny{
		"jti": "123",
		"aud": []string{"t1"},
	}
	assert.Equal(t, `{"aud":["t1"],"jti":"123"}`, c.JSON())
	assert.Equal(t, `aud:
    - t1
jti: "123"
`, c.YAML())
	assert.Equal(t, []string{"t1"}, c.StringSlice("aud"))

	c.GetOrSet("exp", func(key string) any {
		return "123"
	})

	c2 := FromJSON(`{"aud":["t1"],"exp":"123","jti":"123"}`)
	assert.Equal(t, c.JSON(), c2.JSON())
	c2 = FromYAML("aud:\n- t1\njti: \"123\"\nexp: \"123\"\n")
	assert.Equal(t, c.JSON(), c2.JSON())

	var c3 testStruct
	require.NoError(t, c2.To(&c3))
	assert.Equal(t, "123", c3.JTI)
	assert.Equal(t, []string{"t1"}, c3.AUD)

	c4 := FromJSON(`{]`)
	assert.Empty(t, c4)
	assert.Empty(t, c4.StringSlice("aud"))

	c4 = FromYAML(`{]`)
	assert.Empty(t, c4)
	assert.Empty(t, c4.StringSlice("aud"))

	assert.True(t, IsCollection(c4))
	assert.True(t, IsCollection(c["aud"]))
	assert.False(t, IsCollection(c["jti"]))
	assert.False(t, IsCollection(nil))
}

func TestExtract(t *testing.T) {
	js := `{
	"l1s": "string",
	"l1m": {
		"l2s": "string",
		"l2m": {
			"l3s": "string",
			"l3i": 123,
			"l3m": {
				"l4s": "string",
				"l4a": [1,2,3,4]
			}
		}
	},
	"l1a": [1,2,3,4]
}`

	m := FromJSON(js)
	require.NotNil(t, m)
	assert.Nil(t, m.Extract("l1s"))
	assert.Nil(t, m.Extract("l1m", "foo"))
	assert.Nil(t, m.Extract("l1m", "l2m", "l3s"))

	m2 := m.Extract("l1m", "l2m", "l3m")
	require.NotNil(t, m2)
	assert.Equal(t, "string", m2.String("l4s"))
}

type testStruct struct {
	JTI string   `json:"jti"`
	AUD []string `json:"aud"`
}

func TestValues_String(t *testing.T) {
	c := func(o MapAny, k, exp string) {
		act := o.String(k)
		assert.Equal(t, act, exp)
	}

	stru := struct {
		Foo string
		B   bool
		I   int
	}{Foo: "foo", B: true, I: -1}

	o := MapAny{
		"foo":        "bar",
		"blank":      "",
		"count":      uint64(1),
		"strings":    []string{"strings"},
		"empty":      []string{},
		"interfaces": []any{"interfaces"},
		"einterface": []any{},
		"struct":     stru,
		"int":        int(1),
		"int16":      int16(20),
		"int32":      int32(32),
		"int64":      int64(64),
		"uint":       uint(123),
		"uint32":     uint32(132),
		"uint64":     uint64(164),
		"float32":    float32(1.1),
		"float64":    float64(2.64),
		"bool_true":  true,
		"bool_false": false,
	}
	c(o, "foo", "bar")
	c(o, "blank", "")
	c(o, "unknown", "")
	c(o, "count", "1")
	c(o, "strings", "strings")
	c(o, "empty", "")
	c(o, "interfaces", "interfaces")
	c(o, "einterface", "")
	c(o, "struct", `{"Foo":"foo","B":true,"I":-1}`)
	c(o, "int16", "20")
	c(o, "int32", "32")
	c(o, "int64", "64")
	c(o, "uint", "123")
	c(o, "uint32", "132")
	c(o, "uint64", "164")
	c(o, "float32", "1")
	c(o, "float64", "2")
	c(o, "bool_true", "true")
	c(o, "bool_false", "false")
}

func TestValues_JSON(t *testing.T) {
	stru := struct {
		Foo string
		B   bool
		I   int
	}{Foo: "foo", B: true, I: -1}

	o := MapAny{
		"foo":        "bar",
		"blank":      "",
		"count":      uint64(1),
		"ints":       []uint64{1, 2, 3, 4},
		"strings":    []string{"strings"},
		"empty":      []string{},
		"interfaces": []any{"interfaces"},
		"einterface": []any{},
		"struct":     stru,
	}

	js := func(o MapAny, k, exp string) {
		act := JSON(o[k])
		assert.Equal(t, act, exp)
	}
	js(o, "foo", `"bar"`)
	js(o, "blank", `""`)
	js(o, "unknown", "")
	js(o, "count", "1")
	js(o, "ints", `[1,2,3,4]`)
	js(o, "strings", `["strings"]`)
	js(o, "empty", `[]`)
	js(o, "interfaces", `["interfaces"]`)
	js(o, "einterface", "[]")
	js(o, "struct", `{"Foo":"foo","B":true,"I":-1}`)
}

func TestValues_Int(t *testing.T) {
	c := func(o MapAny, k string, exp int) {
		act := o.Int(k)
		assert.Equal(t, exp, act)
	}

	o := MapAny{
		"nil":        nil,
		"struct":     struct{}{},
		"z":          "123",
		"ze":         "abc",
		"n":          int(-1),
		"int":        int(1),
		"int16":      int16(20),
		"int32":      int32(32),
		"int64":      int64(64),
		"uint":       uint(123),
		"uint32":     uint32(132),
		"uint64":     uint64(164),
		"interfaces": []any{1},
		"einterface": []any{},
		"float32":    float32(1.1),
		"float64":    float64(2.64),
	}
	c(o, "nil", 0)
	c(o, "struct", 0)
	c(o, "z", 123)
	c(o, "ze", 0)
	c(o, "n", -1)
	c(o, "int", 1)
	c(o, "int16", 20)
	c(o, "int32", 32)
	c(o, "int64", 64)
	c(o, "uint", 123)
	c(o, "uint32", 132)
	c(o, "uint64", 164)
	c(o, "interfaces", 1)
	c(o, "einterface", 0)
	c(o, "float32", 1)
	c(o, "float64", 2)
}

func TestValues_UInt64(t *testing.T) {
	c := func(o MapAny, k string, exp uint64) {
		act := o.UInt64(k)
		assert.Equal(t, exp, act)
	}

	o := MapAny{
		"nil":        nil,
		"struct":     struct{}{},
		"z":          "123",
		"ze":         "abc",
		"n":          int(-1),
		"int":        int(1),
		"int16":      int16(20),
		"int32":      int32(32),
		"int64":      int64(64),
		"uint":       uint(123),
		"uint32":     uint32(132),
		"uint64":     uint64(164),
		"interfaces": []any{1},
		"einterface": []any{},
		"float32":    float32(1.1),
		"float64":    float64(2.64),
	}
	c(o, "nil", uint64(0))
	c(o, "struct", uint64(0))
	c(o, "z", uint64(123))
	c(o, "ze", uint64(0))
	c(o, "n", uint64(0xffffffffffffffff))
	c(o, "int", uint64(1))
	c(o, "int16", uint64(20))
	c(o, "int32", uint64(32))
	c(o, "int64", uint64(64))
	c(o, "uint", uint64(123))
	c(o, "uint32", uint64(132))
	c(o, "uint64", uint64(164))
	c(o, "interfaces", uint64(1))
	c(o, "einterface", uint64(0))
	c(o, "float32", uint64(1))
	c(o, "float64", uint64(2))
}

func TestValues_Int64(t *testing.T) {
	c := func(o MapAny, k string, exp int64) {
		act := o.Int64(k)
		assert.Equal(t, exp, act)
	}

	o := MapAny{
		"nil":        nil,
		"struct":     struct{}{},
		"z":          "123",
		"ze":         "abc",
		"n":          int(-1),
		"int":        int(1),
		"int16":      int16(20),
		"int32":      int32(32),
		"int64":      int64(64),
		"uint":       uint(123),
		"uint32":     uint32(132),
		"uint64":     uint64(164),
		"interfaces": []any{1},
		"einterface": []any{},
		"float32":    float32(1.1),
		"float64":    float64(2.64),
	}
	c(o, "nil", int64(0))
	c(o, "struct", int64(0))
	c(o, "z", int64(123))
	c(o, "ze", int64(0))
	c(o, "n", int64(-1))
	c(o, "int", int64(1))
	c(o, "int16", int64(20))
	c(o, "int32", int64(32))
	c(o, "int64", int64(64))
	c(o, "uint", int64(123))
	c(o, "uint32", int64(132))
	c(o, "uint64", int64(164))
	c(o, "interfaces", int64(1))
	c(o, "einterface", int64(0))
	c(o, "float32", int64(1))
	c(o, "float64", int64(2))
}

const float64EqualityThreshold = 1e-5

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func TestValues_Float64(t *testing.T) {
	c := func(o MapAny, k string, exp float64) {
		act := o.Float64(k)
		assert.True(t, almostEqual(exp, act), "expected %f, got %f", exp, act)
	}

	o := MapAny{
		"nil":        nil,
		"struct":     struct{}{},
		"z":          "123",
		"ze":         "abc",
		"n":          int(-1),
		"int":        int(1),
		"int16":      int16(20),
		"int32":      int32(32),
		"int64":      int64(64),
		"uint":       uint(123),
		"uint32":     uint32(132),
		"uint64":     uint64(164),
		"interfaces": []any{1},
		"einterface": []any{},
		"float32":    float32(1.1),
		"float64":    float64(2.64),
	}
	c(o, "nil", float64(0))
	c(o, "struct", float64(0))
	c(o, "z", float64(123))
	c(o, "ze", float64(0))
	c(o, "n", float64(-1))
	c(o, "int", float64(1))
	c(o, "int16", float64(20))
	c(o, "int32", float64(32))
	c(o, "int64", float64(64))
	c(o, "uint", float64(123))
	c(o, "uint32", float64(132))
	c(o, "uint64", float64(164))
	c(o, "interfaces", float64(1))
	c(o, "einterface", float64(0))
	c(o, "float32", float64(1.1))
	c(o, "float64", float64(2.64))
}

func TestValues_Float32(t *testing.T) {
	c := func(o MapAny, k string, exp float32) {
		act := o.Float32(k)
		assert.True(t, almostEqual(float64(exp), float64(act)), "expected %f, got %f", exp, act)
	}

	o := MapAny{
		"nil":        nil,
		"struct":     struct{}{},
		"z":          "123",
		"ze":         "abc",
		"n":          int(-1),
		"int":        int(1),
		"int16":      int16(20),
		"int32":      int32(32),
		"int64":      int64(64),
		"uint":       uint(123),
		"uint32":     uint32(132),
		"uint64":     uint64(164),
		"interfaces": []any{1},
		"einterface": []any{},
		"float32":    float32(1.1),
		"float64":    float64(2.64),
	}
	c(o, "nil", float32(0))
	c(o, "struct", float32(0))
	c(o, "z", float32(123))
	c(o, "ze", float32(0))
	c(o, "n", float32(-1))
	c(o, "int", float32(1))
	c(o, "int16", float32(20))
	c(o, "int32", float32(32))
	c(o, "int64", float32(64))
	c(o, "uint", float32(123))
	c(o, "uint32", float32(132))
	c(o, "uint64", float32(164))
	c(o, "interfaces", float32(1))
	c(o, "einterface", float32(0))
	c(o, "float32", float32(1.1))
	c(o, "float64", float32(2.64))
}

func TestValues_Bool(t *testing.T) {
	c := func(o MapAny, k string, exp bool) {
		act := o.Bool(k)
		assert.Equal(t, act, exp)
	}

	o := MapAny{
		"nil":        nil,
		"struct":     struct{}{},
		"true":       true,
		"false":      false,
		"strue":      "true",
		"interfaces": []any{"true"},
		"einterface": []any{},
	}
	c(o, "nil", false)
	c(o, "struct", false)
	c(o, "true", true)
	c(o, "strue", true)
	c(o, "false", false)
	c(o, "interfaces", true)
	c(o, "einterface", false)
}

func TestValues_Time(t *testing.T) {
	c := func(o MapAny, k string, exp *time.Time) {
		t.Run(k, func(t *testing.T) {
			act := o.Time(k)
			assert.Equal(t, exp, act, k)
		})
	}

	tPtr := func(val time.Time) *time.Time {
		return &val
	}
	loc := time.FixedZone("", -25200)
	o := MapAny{
		"nil":        nil,
		"struct":     struct{}{},
		"z":          "123",
		"ze":         "abc",
		"time":       time.Unix(123, 0),
		"*time":      tPtr(time.Unix(123, 0)),
		"str":        "2006-01-02T15:04:05.000-0700",
		"invalid":    "2006000000000000000000000000",
		"int":        int(189898989898),
		"float":      float64(189898989898.0),
		"int16":      int16(20),
		"int32":      int32(32),
		"int64":      int64(64),
		"uint":       uint(123),
		"uint32":     uint32(132),
		"uint64":     uint64(164),
		"interfaces": []any{uint64(164)},
		"einterface": []any{},
	}
	c(o, "nil", nil)
	c(o, "struct", nil)
	c(o, "z", tPtr(time.Unix(123, 0)))
	c(o, "time", tPtr(time.Unix(123, 0)))
	c(o, "*time", tPtr(time.Unix(123, 0)))
	c(o, "str", tPtr(time.Date(2006, time.January, 2, 15, 4, 5, 0, loc)))
	c(o, "invalid", nil)
	c(o, "ze", nil)
	c(o, "int16", nil)
	c(o, "int32", nil)
	c(o, "int", tPtr(time.Unix(189898989898, 0)))
	c(o, "float", tPtr(time.Unix(189898989898, 0)))
	c(o, "int64", tPtr(time.Unix(64, 0)))
	c(o, "uint", nil)
	c(o, "uint32", nil)
	c(o, "uint64", tPtr(time.Unix(164, 0)))
	c(o, "interfaces", tPtr(time.Unix(164, 0)))
	c(o, "einterface", nil)
}

func TestValues_Empty(t *testing.T) {
	var none MapAny
	assert.Equal(t, "", none.String("1"))
	assert.Equal(t, false, none.Bool("1"))
	assert.Nil(t, none.Time("1"))
	assert.Equal(t, 0, none.Int("1"))
	assert.Equal(t, uint64(0), none.UInt64("1"))
	assert.Equal(t, int64(0), none.Int64("1"))
}

func TestValues_StringSlice(t *testing.T) {
	r := StringSlice(nil)
	assert.Equal(t, []string{}, r)

	r = StringSlice([]string{"str1", "str2", "str3"})
	assert.Len(t, r, 3)
	assert.Equal(t, []string{"str1", "str2", "str3"}, r)

	r = StringSlice([]any{1, "str2", 5})
	assert.Len(t, r, 3)
	assert.Equal(t, []string{"1", "str2", "5"}, r)

	r = StringSlice([]any{"str1", "str2", "str3"})
	assert.Len(t, r, 3)
	assert.Equal(t, []string{"str1", "str2", "str3"}, r)

	r = StringSlice([]int{1, 2, 3})
	assert.Equal(t, []string{}, r)
}

func Test_NvlInt(t *testing.T) {
	c := func(exp int, items ...int) {
		act := NumbersCoalesce(items...)
		if act != exp {
			t.Errorf("Expecting NvlInt(%v) to return %d, but got %d", items, exp, act)
		}
	}
	c(0)
	c(0, 0)
	c(10, 10)
	c(10, 10, 0)
	c(-10, -10)
	c(10, 0, 10)
	c(-5, 0, -5, 10)
}

func Test_NvlInt64(t *testing.T) {
	c := func(exp int64, items ...int64) {
		act := NumbersCoalesce(items...)
		if act != exp {
			t.Errorf("Expecting NvlInt64(%v) to return %d, but got %d", items, exp, act)
		}
	}
	c(0)
	c(0, 0)
	c(10, 10)
	c(10, 10, 0)
	c(-10, -10)
	c(10, 0, 10)
	c(-5, 0, -5, 10)
}

func Test_NvlUint64(t *testing.T) {
	c := func(exp uint64, items ...uint64) {
		act := NumbersCoalesce(items...)
		if act != exp {
			t.Errorf("Expecting NvlUnt64(%v) to return %d, but got %d", items, exp, act)
		}
	}
	c(0)
	c(0, 0)
	c(10, 10)
	c(10, 10, 0)
	c(10, 0, 10)
	c(5, 0, 5, 10)
	c(5, 0, 5, 0)
}

func Test_StringsCoalesce(t *testing.T) {
	assert.Equal(t, "", StringsCoalesce())
	assert.Equal(t, "1", StringsCoalesce("1", "2", "3"))
	assert.Equal(t, "2", StringsCoalesce("", "2", "3"))
	assert.Equal(t, "3", StringsCoalesce("", "", "3"))
}

func Test_Coalesce(t *testing.T) {
	assert.Equal(t, "1", Coalesce("1", "2", "3"))
	assert.Equal(t, "2", Coalesce("", "2", "3"))
	assert.Equal(t, "3", Coalesce("", "", "3"))
	assert.Equal(t, "", Coalesce(""))

	assert.Equal(t, []string{"1"}, Coalesce([]string{"1"}, []string{"2", "3"}))
	assert.Equal(t, []string{""}, Coalesce([]string{""}, []string{"2", "3"}))
	assert.Equal(t, []string{"2", "3"}, Coalesce([]string{}, []string{"2", "3"}))
	var empty []string
	assert.Equal(t, []string{"3"}, Coalesce(empty, empty, []string{"3"}))
}

func TestNvlNumber(t *testing.T) {
	assert.Equal(t, 1, NumbersCoalesce(0, 1))
	assert.Equal(t, uint64(1), NumbersCoalesce(0, uint64(1)))
}

func TestSelect(t *testing.T) {
	assert.Equal(t, 1, Select(false, 0, 1))
	assert.Equal(t, uint64(0), Select(true, 0, uint64(1)))
}

func TestRangeOrderedMap(t *testing.T) {
	m1 := map[string]string{"ya": "1", "b": "2", "c": "3"}

	var keys []string
	RangeOrderedMap(m1, func(k, v string) bool {
		keys = append(keys, k)
		return true
	})
	assert.Equal(t, []string{"b", "c", "ya"}, keys)
}

func Test_MapAny_DB(t *testing.T) {
	tcases := []struct {
		val MapAny
		exp string
	}{
		{val: MapAny{"one": "two"}, exp: "{\"one\":\"two\"}"},
		{val: MapAny{}, exp: ""},
		{val: nil, exp: ""},
	}

	for _, tc := range tcases {
		dr, err := tc.val.Value()
		require.NoError(t, err)

		var drv string
		if v, ok := dr.(string); ok {
			drv = v
		}
		assert.Equal(t, tc.exp, drv)

		var val2 MapAny
		err = val2.Scan(dr)
		require.NoError(t, err)
		assert.Equal(t, len(tc.val), len(val2))
	}

	m1 := MapAny{"one": "two"}
	m2 := MapAny{"three": "four"}

	m1.Merge(m2)
	assert.Equal(t, 2, len(m1))

	var mm MapAny
	mm.Merge(m2)
	assert.Equal(t, 1, len(mm))
	mm.Merge(m1)
	assert.Equal(t, 2, len(mm))
}
