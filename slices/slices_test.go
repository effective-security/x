package slices

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlices_NvlString(t *testing.T) {
	v := func(exp string, items ...string) {
		act := NvlString(items...)
		if act != exp {
			t.Errorf("Expecting NvlString(%v) to return %s, but got %s", items, exp, act)
		}
	}
	v("")
	v("", "")
	v("", "", "")
	v("a", "a")
	v("a", "a", "")
	v("a", "", "", "a")
	v("b", "", "b", "a")
}

func TestSlices_CloneStrings(t *testing.T) {
	c := CloneStrings(nil)
	if c != nil {
		t.Errorf("CloneStrings with a nil src should return nil, but returned %#v", c)
	}
	c = CloneStrings([]string{})
	if c == nil || len(c) != 0 {
		t.Errorf("CloneStrings with a non-nil, but empty slice, should return a new empty slice, but got %#v", c)
	}
	s := []string{"a", "b", "c"}
	c = CloneStrings(s)
	if !reflect.DeepEqual(s, c) {
		t.Errorf("CloneString() returned different contents to source, got %+v, expecting %+v", c, s)
	}
	s[0] = "x"
	if c[0] != "a" {
		t.Errorf("CloneString() didn't return a Clone, it was mutated by mutating the source")
	}
}

func TestSlices_ContainsString(t *testing.T) {
	s := []string{"a", "b", "c", "foo", "bar", "qux"}
	missing := []string{"bob", "quxx"}
	testSlicesContains(t, s, missing, "q", func(items any, item any) bool {
		return ContainsString(items.([]string), item.(string))
	})
	if ContainsString(nil, "") {
		t.Errorf("a nil slice shouldn't contain anything!")
	}
}

func TestSlices_ContainsStringEqualFold(t *testing.T) {
	src := []string{"one", "TWO", "Three"}
	tests := []string{"ONE", "One", "two", "three"}
	m := []string{"", "oned", "Four"}
	for _, item := range append(src, tests...) {
		if !ContainsStringEqualFold(src, item) {
			t.Errorf("Expecting to find %q in %v, but didn't", item, src)
		}
	}
	for _, item := range m {
		if ContainsStringEqualFold(src, item) {
			t.Errorf("Not expecting to find %q in %v, but did", item, src)
		}
	}
}

func TestSlices_StringContainsOneOf(t *testing.T) {
	tcases := []struct {
		str    string
		slices []string
		exp    bool
	}{
		{"Daniel", []string{"foo", "bar"}, false},
		{"Daniel", []string{"foo", "el"}, true},
		{"Daniel", []string{"foo", "da"}, false},
		{"Daniel", []string{"foo", "Dan"}, true},
	}
	for idx, tc := range tcases {
		res := StringContainsOneOf(tc.str, tc.slices)
		if res != tc.exp {
			t.Errorf("case %d failed", idx)
		}
	}
}

func TestSlices_StringStartsWithOneOf(t *testing.T) {
	tcases := []struct {
		str    string
		slices []string
		exp    bool
	}{
		{"Daniel", []string{"foo", "bar"}, false},
		{"foo_Daniel", []string{"foo", "el"}, true},
		{"daniel", []string{"foo", "da"}, true},
		{"Daniel", []string{"foo", "Dan"}, true},
	}
	for idx, tc := range tcases {
		res := StringStartsWithOneOf(tc.str, tc.slices)
		if res != tc.exp {
			t.Errorf("case %d failed", idx)
		}
	}
}

func testSlicesContains(t *testing.T, items any, missing any, newItem any, containsFunc func(items any, item any) bool) {
	vm := reflect.ValueOf(missing)
	for i := 0; i < vm.Len(); i++ {
		if containsFunc(items, vm.Index(i).Interface()) {
			t.Errorf("Item %v wasn't in items slice, but contains said it was!", vm.Index(i))
		}
	}
	vi := reflect.ValueOf(items)
	for i := 0; i < vi.Len(); i++ {
		if !containsFunc(items, vi.Index(i).Interface()) {
			t.Errorf("Item %v is at index %d in slice, but contains said it wasn't in the slice", vi.Index(i), i)
		}
	}
	vi = reflect.Append(vi, reflect.ValueOf(newItem))
	if !containsFunc(vi.Interface(), newItem) {
		t.Errorf("Item %v was added to slice, but contains didn't spot it", newItem)
	}
	if containsFunc(vi.Slice(1, vi.Len()-1).Interface(), vi.Index(0).Interface()) {
		t.Errorf("Item %v wasn't in the modified slice, but contains said it was", vi.Index(0))
	}
}

func TestSlices_ByteSlicesEqual(t *testing.T) {
	bytes := []any{
		[]byte{},
		[]byte{1},
		[]byte{1, 2, 3},
		[]byte{1, 2, 3, 4},
		[]byte{2, 2, 3, 4},
		[]byte{1, 2, 3, 5},
	}
	testSlicesEquals(t, "Byte", bytes, bytes[2], []byte{1, 2, 3}, func(x, y any) bool {
		return ByteSlicesEqual(x.([]byte), y.([]byte))
	})
	if ByteSlicesEqual(nil, []byte{1}) || ByteSlicesEqual([]byte{1}, nil) {
		t.Errorf("ByteSliceEqual for a nil slice shouldn't return true when the other slice has items in it")
	}
	if !ByteSlicesEqual(nil, nil) || !ByteSlicesEqual(nil, []byte{}) {
		t.Errorf("ByteSlicesEquals for a nil & empty slice should return true")
	}
}

func TestSlices_StringSlicesEqual(t *testing.T) {
	strings := []any{
		[]string{},
		[]string{""},
		[]string{"aa"},
		[]string{"aa", "bb"},
		[]string{"aa", "bb", "cc"},
		[]string{"bb", "bb", "cc"},
		[]string{"aa", "bb", "bb"},
	}
	testSlicesEquals(t, "String", strings, []string{"aa", "bb", "cc"}, strings[4], func(x, y any) bool {
		return StringSlicesEqual(x.([]string), y.([]string))
	})
	if StringSlicesEqual(nil, []string{"a"}) || StringSlicesEqual([]string{"a"}, nil) {
		t.Errorf("StringSlicesEqual for nil and a slice with an item in it should return false")
	}
	if !StringSlicesEqual(nil, nil) || !StringSlicesEqual(nil, []string{}) {
		t.Errorf("StringSlicesEqual for a nil and empty slice should return true")
	}
}

func assertStringSlicesEqual(t *testing.T, preamble string, exp []string, act []string) {
	if len(act) != len(exp) {
		t.Errorf("%s: expected to get %d items, but got %d", preamble, len(exp), len(act))
	} else {
		for i, a := range act {
			if a != exp[i] {
				t.Errorf("%s: at index %d expected to get %q, but got %q", preamble, i, exp[i], a)
			}
		}
	}
}

func TestSlices_Quoted(t *testing.T) {
	c := func(in, exp []string) {
		res := Quoted(in)
		assertStringSlicesEqual(t, fmt.Sprintf("Quoted(%v)", in), exp, res)
	}
	c([]string{}, []string{})
	c([]string{"bob "}, []string{`"bob "`})
	c([]string{"b", "a", "c"}, []string{`"b"`, `"a"`, `"c"`})
}

func TestSlices_Prefixed(t *testing.T) {
	c := func(p string, items []string, exp []string) {
		act := Prefixed(p, items)
		assertStringSlicesEqual(t, fmt.Sprintf("Prefixed(%v,%v)", p, items), exp, act)
	}
	c("bob", []string{}, []string{})
	c("bob", []string{"alice"}, []string{"bobalice"})
	c("bob", []string{"alice", "eve"}, []string{"bobalice", "bobeve"})
	c("", []string{"alice", "eve"}, []string{"alice", "eve"})
}

func TestSlices_Suffix(t *testing.T) {
	c := func(p string, items []string, exp []string) {
		act := Suffixed(p, items)
		assertStringSlicesEqual(t, fmt.Sprintf("Suffixed(%v,%v)", p, items), exp, act)
	}
	c("bob", []string{}, []string{})
	c("bob", []string{"alice"}, []string{"alicebob"})
	c("bob", []string{"alice", "eve"}, []string{"alicebob", "evebob"})
	c("", []string{"alice", "eve"}, []string{"alice", "eve"})
}

func TestSlices_Int64SlicesEqual(t *testing.T) {
	vals := []any{
		[]int64{},
		[]int64{0},
		[]int64{1},
		[]int64{42, 43},
		[]int64{42, 43, 0},
		[]int64{41, 43, 0},
		[]int64{42, 43, 43},
	}
	testSlicesEquals(t, "Int64", vals, []int64{42, 43, 0}, vals[4], func(x, y any) bool {
		return Int64SlicesEqual(x.([]int64), y.([]int64))
	})
	if Int64SlicesEqual(nil, []int64{1}) || Int64SlicesEqual([]int64{1}, nil) {
		t.Errorf("Int64SlicesEqual for a nil slice and a slice with items should return false")
	}
	if !Int64SlicesEqual(nil, nil) || !Int64SlicesEqual(nil, []int64{}) {
		t.Errorf("Int64SlicesEqual for a nil slice and an empty slice should return true")
	}
}

func TestSlices_UInt64SlicesEqual(t *testing.T) {
	vals := []any{
		[]uint64{},
		[]uint64{0},
		[]uint64{1},
		[]uint64{42, 43},
		[]uint64{42, 43, 0},
		[]uint64{41, 43, 0},
		[]uint64{42, 43, 43},
	}
	testSlicesEquals(t, "Uint64", vals, []uint64{42, 43, 0}, vals[4], func(x, y any) bool {
		return Uint64SlicesEqual(x.([]uint64), y.([]uint64))
	})
	if Uint64SlicesEqual(nil, []uint64{1}) || Uint64SlicesEqual([]uint64{1}, nil) {
		t.Errorf("Uint64SlicesEqual for a nil slice and a slice with items should return false")
	}
	if !Uint64SlicesEqual(nil, nil) || !Uint64SlicesEqual(nil, []uint64{}) {
		t.Errorf("Uint64SlicesEqual for a nil slice and an empty slice should return true")
	}
}

func TestSlices_BoolSlicesEqual(t *testing.T) {
	bools := []any{
		[]bool{},
		[]bool{false},
		[]bool{true},
		[]bool{false, false},
		[]bool{false, false, true},
		[]bool{true, false, true},
		[]bool{false, false, false},
	}
	testSlicesEquals(t, "Bool", bools, []bool{false, false, true}, bools[4], func(x, y any) bool {
		return BoolSlicesEqual(x.([]bool), y.([]bool))
	})
	if BoolSlicesEqual(nil, []bool{false}) || BoolSlicesEqual([]bool{false}, nil) {
		t.Errorf("BoolSlicesEqual for a nil and slice with items should return false")
	}
	if !BoolSlicesEqual(nil, nil) || !BoolSlicesEqual(nil, []bool{}) {
		t.Errorf("BoolSlicesEqual for a nil and empty slice should return true")
	}
}

func TestSlices_FloatSlicesEqual(t *testing.T) {
	vals := []any{
		[]float64{},
		[]float64{0},
		[]float64{1, 2},
		[]float64{3, 4, 5},
		[]float64{2.0, 4, 5},
		[]float64{3, 4, 4},
	}
	testSlicesEquals(t, "Float64", vals, []float64{2.0, 4, 5}, vals[4], func(x, y any) bool {
		return Float64SlicesEqual(x.([]float64), y.([]float64))
	})
	if Float64SlicesEqual(nil, []float64{0}) || Float64SlicesEqual([]float64{0}, nil) {
		t.Errorf("Float64SlicesEqual for a nil and slice with items should return false")
	}
	if !Float64SlicesEqual(nil, nil) || !Float64SlicesEqual(nil, []float64{}) {
		t.Errorf("Float64SlicesEqual for a nil and empty slice should return true")
	}
}

func testSlicesEquals(t *testing.T, funcName string, vals []any, goodVal1 any, goodVal2 any, equalsFunc func(x, y any) bool) {
	for i, x := range vals {
		for j, y := range vals {
			r := equalsFunc(x, y)
			if (i == j) && !r {
				t.Errorf("%vSlicesEqual for the same slice shouldn't return false! (%v,%v)", funcName, x, y)
			} else if (i != j) && r {
				t.Errorf("%vSlicesEqual for different slices should return false! (%v,%v)", funcName, x, y)
			}
		}
	}
	if !equalsFunc(goodVal1, goodVal2) {
		t.Errorf("Different slices with the same contents should return true for %vSlicesEqual (%v,%v)", funcName, goodVal1, goodVal2)
	}
}

func TestSlices_StringUpto(t *testing.T) {
	assert.Equal(t, "", StringUpto("", 0))
	assert.Equal(t, "", StringUpto("", 2))
	assert.Equal(t, "", StringUpto("11", 0))
	assert.Equal(t, "1", StringUpto("11", 1))
	assert.Equal(t, "11", StringUpto("11", 2))
	assert.Equal(t, "11", StringUpto("11", 3))
}

func Test_removeDuplicates(t *testing.T) {
	dups := []string{"12", "45", "45", "78", "12", "porto"}
	noDups := UniqueStrings(dups)
	assert.Equal(t, len(noDups), 4)
}

func Test_Replace(t *testing.T) {
	Replace(nil, "", "")

	sliceStr := []string{}
	Replace(sliceStr, "old", "new")
	assert.Equal(t, 0, len(sliceStr))

	sliceStr = []string{"old", "other"}
	Replace(sliceStr, "", "something")
	assert.Equal(t, []string{"old", "other"}, sliceStr)

	sliceStr = []string{"old", "other"}
	Replace(sliceStr, "old", "new")
	assert.Equal(t, []string{"new", "other"}, sliceStr)

	sliceInt := []int{1}
	Replace(sliceInt, 2, 4)
	assert.Equal(t, []int{1}, sliceInt)

	sliceInt = []int{1, 2, 3}
	Replace(sliceInt, 2, 4)
	assert.Equal(t, []int{1, 4, 3}, sliceInt)
}

func Test_StringArrayToMap(t *testing.T) {
	m, err := StringArrayToMap(nil)
	require.NoError(t, err)
	assert.Equal(t, 0, len(m))

	m, err = StringArrayToMap([]string{})
	require.NoError(t, err)
	assert.Equal(t, 0, len(m))

	m, err = StringArrayToMap([]string{"a="}) // value is optional
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"a": ""}, m)

	m, err = StringArrayToMap([]string{"a=b", "c=d"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "b", "c": "d"}, m)

	m, err = StringArrayToMap([]string{"a1,a2,a3=b1,b2,b3", "c1,c2,c3=d1,d2,d3"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"a1,a2,a3": "b1,b2,b3", "c1,c2,c3": "d1,d2,d3"}, m)

	m, err = StringArrayToMap([]string{"a=b1=b2=b3", "c=d1=d2=d3"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "b1=b2=b3", "c": "d1=d2=d3"}, m)

	// error cases
	m, err = StringArrayToMap([]string{"=a"}) // no key
	require.Error(t, err)
	assert.Equal(t, 0, len(m))

	m, err = StringArrayToMap([]string{"a"}) // no separator
	require.Error(t, err)
	assert.Equal(t, 0, len(m))
}

func Test_TruncateArray(t *testing.T) {
	assert.Nil(t, Truncate([]any(nil), 0))
	assert.Nil(t, Truncate([]any(nil), 5))
	assert.Empty(t, Truncate([]string{}, 0))
	assert.Empty(t, Truncate([]string{}, 5))
	assert.True(t, reflect.DeepEqual([]int{}, Truncate([]int{1, 2, 3}, 0)))
	assert.True(t, reflect.DeepEqual([]int{1}, Truncate([]int{1, 2, 3}, 1)))
	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, Truncate([]int{1, 2, 3}, 3)))
}

func Test_Contains(t *testing.T) {
	assert.False(t, Contains([]int(nil), 1))
	assert.False(t, Contains([]int{}, 1))
	assert.False(t, Contains([]int{2, 3, 4}, 1))
	assert.True(t, Contains([]int{1}, 1))
	assert.True(t, Contains([]int{2, 3, 4}, 3))

	assert.False(t, Contains([]string{""}, "a"))
	assert.True(t, Contains([]string{""}, ""))
	assert.True(t, Contains([]string{"a", "b", ""}, ""))
}

func Test_Hash(t *testing.T) {
	require.NotEmpty(t, HashStrings())
	require.NotEmpty(t, HashStrings("a"))
	require.NotEmpty(t, HashStrings("a", "b", "c"))
}

func Test_Deduplicate(t *testing.T) {
	assert.Equal(t, 0, len(Deduplicate[int](nil)))
	assert.Equal(t, 0, len(Deduplicate([]int{})))
	assert.Equal(t, []int{1}, Deduplicate([]int{1}))
	assert.Equal(t, []int{1, 2, 3}, Deduplicate([]int{1, 2, 2, 3, 1, 1}))
}
