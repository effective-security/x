package maps_test

import (
	stdmaps "maps"
	"slices"
	"sort"
	"strconv"
	"testing"

	xmaps "github.com/effective-security/x/maps"
	"github.com/stretchr/testify/assert"
)

const (
	smallMapSize = 16
	largeMapSize = 1024
)

var (
	benchmarkMapKeys   []int
	benchmarkMapValues []int
	benchmarkMapResult map[int]int
)

func TestStdlibMapCompatibility(t *testing.T) {
	input := map[string]int{
		"a": 3,
		"c": 1,
		"z": 2,
	}
	assert.ElementsMatch(t, slices.Collect(stdmaps.Keys(input)), xmaps.Keys(input))
	assert.ElementsMatch(t, slices.Collect(stdmaps.Values(input)), xmaps.Values(input))
	assert.Equal(t, slices.Sorted(stdmaps.Keys(input)), xmaps.OrderedKeys(input))

	left := map[string]int{
		"a": 1,
		"b": 2,
	}
	right := map[string]int{
		"b": 3,
		"c": 4,
	}
	assert.Equal(t, mergeStdlib(left, right), xmaps.Merge(left, right))

	localEmpty := xmaps.Keys(map[string]int{})
	stdlibEmpty := slices.Collect(stdmaps.Keys(map[string]int{}))
	assert.NotNil(t, localEmpty)
	assert.Nil(t, stdlibEmpty)
}

func BenchmarkKeysVsStdlib(b *testing.B) {
	for _, size := range []int{smallMapSize, largeMapSize} {
		input := benchmarkInputMap(size)
		b.Run(strconv.Itoa(size)+"/Local", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkMapKeys = xmaps.Keys(input)
			}
		})
		b.Run(strconv.Itoa(size)+"/Stdlib", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkMapKeys = slices.Collect(stdmaps.Keys(input))
			}
		})
		b.Run(strconv.Itoa(size)+"/StdlibPreallocated", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				keys := make([]int, 0, len(input))
				benchmarkMapKeys = slices.AppendSeq(keys, stdmaps.Keys(input))
			}
		})
	}
}

func BenchmarkValuesVsStdlib(b *testing.B) {
	for _, size := range []int{smallMapSize, largeMapSize} {
		input := benchmarkInputMap(size)
		b.Run(strconv.Itoa(size)+"/Local", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkMapValues = xmaps.Values(input)
			}
		})
		b.Run(strconv.Itoa(size)+"/Stdlib", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkMapValues = slices.Collect(stdmaps.Values(input))
			}
		})
		b.Run(strconv.Itoa(size)+"/StdlibPreallocated", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				values := make([]int, 0, len(input))
				benchmarkMapValues = slices.AppendSeq(values, stdmaps.Values(input))
			}
		})
	}
}

func BenchmarkOrderedKeysVsStdlib(b *testing.B) {
	input := benchmarkInputMap(largeMapSize)
	b.Run("LegacySortSlice", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkMapKeys = legacyOrderedKeys(input)
		}
	})
	b.Run("LocalSlicesSort", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkMapKeys = xmaps.OrderedKeys(input)
		}
	})
	b.Run("Stdlib", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkMapKeys = slices.Sorted(stdmaps.Keys(input))
		}
	})
}

func BenchmarkMergeVsStdlib(b *testing.B) {
	left := benchmarkInputMap(largeMapSize)
	right := benchmarkInputMap(largeMapSize * 2)
	b.Run("LegacyLoops", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkMapResult = legacyMerge(left, right)
		}
	})
	b.Run("CurrentWrapper", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkMapResult = xmaps.Merge(left, right)
		}
	})
	b.Run("StdlibCloneCopy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkMapResult = mergeStdlib(left, right)
		}
	})
}

func benchmarkInputMap(size int) map[int]int {
	result := make(map[int]int, size)
	for i := 0; i < size; i++ {
		result[i] = i
	}
	return result
}

func legacyOrderedKeys(input map[int]int) []int {
	keys := xmaps.Keys(input)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func mergeStdlib[K comparable, V any](inputs ...map[K]V) map[K]V {
	if len(inputs) == 0 {
		return make(map[K]V)
	}
	result := stdmaps.Clone(inputs[0])
	for _, input := range inputs[1:] {
		stdmaps.Copy(result, input)
	}
	return result
}

func legacyMerge[K comparable, V any](inputs ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, input := range inputs {
		for key, value := range input {
			result[key] = value
		}
	}
	return result
}
