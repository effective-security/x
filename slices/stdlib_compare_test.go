package slices_test

import (
	"cmp"
	stdslices "slices"
	"sort"
	"testing"

	xslices "github.com/effective-security/x/slices"
	"github.com/stretchr/testify/assert"
)

const benchmarkSliceSize = 1024

var (
	benchmarkBoolResult  bool
	benchmarkSliceResult []string
	benchmarkString      string
	benchmarkUint64s     []uint64
)

func TestStdlibSliceCompatibility(t *testing.T) {
	strings := []string{"a", "b", "c"}
	assert.Equal(t, stdslices.Contains(strings, "b"), xslices.ContainsString(strings, "b"))
	assert.Equal(t, stdslices.Contains(strings, "b"), xslices.Contains(strings, "b"))
	assert.Equal(t, stdslices.Equal(strings, stdslices.Clone(strings)), xslices.StringSlicesEqual(strings, xslices.CloneStrings(strings)))
	assert.Equal(t, cmp.Or("", "value"), xslices.NvlString("", "value"))
	assert.Nil(t, xslices.CloneStrings(nil))
	assert.Nil(t, stdslices.Clone([]string(nil)))
}

func BenchmarkContainsVsStdlib(b *testing.B) {
	input := make([]int, benchmarkSliceSize)
	missing := benchmarkSliceSize + 1
	b.Run("LegacyLoop", func(b *testing.B) {
		for b.Loop() {
			benchmarkBoolResult = legacyContains(input, missing)
		}
	})
	b.Run("CurrentWrapper", func(b *testing.B) {
		for b.Loop() {
			benchmarkBoolResult = xslices.Contains(input, missing)
		}
	})
	b.Run("Stdlib", func(b *testing.B) {
		for b.Loop() {
			benchmarkBoolResult = stdslices.Contains(input, missing)
		}
	})
}

func BenchmarkEqualVsStdlib(b *testing.B) {
	left := make([]string, benchmarkSliceSize)
	right := stdslices.Clone(left)
	b.Run("LegacyLoop", func(b *testing.B) {
		for b.Loop() {
			benchmarkBoolResult = legacyStringsEqual(left, right)
		}
	})
	b.Run("CurrentWrapper", func(b *testing.B) {
		for b.Loop() {
			benchmarkBoolResult = xslices.StringSlicesEqual(left, right)
		}
	})
	b.Run("Stdlib", func(b *testing.B) {
		for b.Loop() {
			benchmarkBoolResult = stdslices.Equal(left, right)
		}
	})
}

func BenchmarkCloneVsStdlib(b *testing.B) {
	input := make([]string, benchmarkSliceSize)
	b.Run("LegacyMakeCopy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result := make([]string, len(input))
			copy(result, input)
			benchmarkSliceResult = result
		}
	})
	b.Run("CurrentWrapper", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSliceResult = xslices.CloneStrings(input)
		}
	})
	b.Run("Stdlib", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSliceResult = stdslices.Clone(input)
		}
	})
}

func BenchmarkNvlStringVsStdlib(b *testing.B) {
	values := []string{"", "", "value", "unused"}
	b.Run("LegacyLoop", func(b *testing.B) {
		for b.Loop() {
			benchmarkString = legacyNvlString(values...)
		}
	})
	b.Run("CurrentWrapper", func(b *testing.B) {
		for b.Loop() {
			benchmarkString = xslices.NvlString(values...)
		}
	})
	b.Run("Stdlib", func(b *testing.B) {
		for b.Loop() {
			benchmarkString = cmp.Or(values...)
		}
	})
}

func BenchmarkUint64SortVsStdlib(b *testing.B) {
	input := make([]uint64, benchmarkSliceSize)
	for i := range input {
		input[i] = uint64(benchmarkSliceSize - i)
	}
	b.Run("LocalSortInterface", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result := stdslices.Clone(input)
			sort.Sort(xslices.Uint64s(result))
			benchmarkUint64s = result
		}
	})
	b.Run("Stdlib", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result := stdslices.Clone(input)
			stdslices.Sort(result)
			benchmarkUint64s = result
		}
	})
}

func legacyContains[T comparable](values []T, target T) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func legacyStringsEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i, value := range left {
		if value != right[i] {
			return false
		}
	}
	return true
}

func legacyNvlString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
