package maps_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/effective-security/x/maps"
)

func TestKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []string
	}{
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: []string{},
		},
		{
			name:     "single element",
			input:    map[string]int{"a": 1},
			expected: []string{"a"},
		},
		{
			name:     "multiple elements",
			input:    map[string]int{"a": 1, "b": 2, "c": 3},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maps.Keys(tt.input)
			// Sort both slices for comparison since map iteration order is not guaranteed
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []int
	}{
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: []int{},
		},
		{
			name:     "single element",
			input:    map[string]int{"a": 1},
			expected: []int{1},
		},
		{
			name:     "multiple elements",
			input:    map[string]int{"a": 1, "b": 2, "c": 3},
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maps.Values(tt.input)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestEach(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3}
	visited := make(map[string]int)

	maps.Each(input, func(key string, value int) {
		visited[key] = value
	})

	assert.Equal(t, input, visited)
}

func TestMap(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test value transformation
	result := maps.Map(input, func(key string, value int) string {
		return fmt.Sprintf("%s:%d", key, value)
	})

	expected := map[string]string{
		"a": "a:1",
		"b": "b:2",
		"c": "c:3",
	}
	assert.Equal(t, expected, result)

	// Test key-value transformation
	result2 := maps.Map(input, func(key string, value int) int {
		return len(key) + value
	})

	expected2 := map[string]int{
		"a": 2, // len("a") + 1
		"b": 3, // len("b") + 2
		"c": 4, // len("c") + 3
	}
	assert.Equal(t, expected2, result2)
}

func TestFilter(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}

	// Filter even values
	result := maps.Filter(input, func(key string, value int) bool {
		return value%2 == 0
	})

	expected := map[string]int{"b": 2, "d": 4}
	assert.Equal(t, expected, result)

	// Filter keys starting with 'a'
	result2 := maps.Filter(input, func(key string, value int) bool {
		return key[0] == 'a'
	})

	expected2 := map[string]int{"a": 1}
	assert.Equal(t, expected2, result2)
}

func TestReduce(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3}

	// Sum all values
	result := maps.Reduce(input, 0, func(acc int, key string, value int) int {
		return acc + value
	})

	assert.Equal(t, 6, result)

	// Concatenate keys
	result2 := maps.Reduce(input, "", func(acc string, key string, value int) string {
		return acc + key
	})

	// Order might vary, so check length and that all keys are present
	assert.Equal(t, 3, len(result2))
	assert.Contains(t, result2, "a")
	assert.Contains(t, result2, "b")
	assert.Contains(t, result2, "c")
}

func TestFind(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	// Find first value > 2
	key, value, found := maps.Find(input, func(key string, value int) bool {
		return value > 2
	})

	assert.True(t, found)
	assert.Contains(t, []string{"c", "d"}, key)
	assert.Contains(t, []int{3, 4}, value)

	// Find non-existent
	_, _, found = maps.Find(input, func(key string, value int) bool {
		return value > 10
	})

	assert.False(t, found)
}

func TestAny(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test with existing match
	result := maps.Any(input, func(key string, value int) bool {
		return value > 2
	})
	assert.True(t, result)

	// Test with no match
	result = maps.Any(input, func(key string, value int) bool {
		return value > 10
	})
	assert.False(t, result)
}

func TestAll(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test with all matching
	result := maps.All(input, func(key string, value int) bool {
		return value > 0
	})
	assert.True(t, result)

	// Test with some not matching
	result = maps.All(input, func(key string, value int) bool {
		return value > 1
	})
	assert.False(t, result)
}

func TestMerge(t *testing.T) {
	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"b": 3, "c": 4}
	map3 := map[string]int{"d": 5}

	// Merge two maps
	result := maps.Merge(map1, map2)
	expected := map[string]int{"a": 1, "b": 3, "c": 4}
	assert.Equal(t, expected, result)

	// Merge three maps
	result = maps.Merge(map1, map2, map3)
	expected = map[string]int{"a": 1, "b": 3, "c": 4, "d": 5}
	assert.Equal(t, expected, result)

	// Merge empty
	result = maps.Merge[string, int]()
	assert.Equal(t, map[string]int{}, result)
}

func TestInvert(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3, "d": 1, "e": 2}

	result := maps.Invert(input)
	expected := map[int][]string{1: {"a", "d"}, 2: {"b", "e"}, 3: {"c"}}
	require.Equal(t, len(expected), len(result))
	for k, v := range expected {
		assert.ElementsMatch(t, v, result[k])
	}
}

func TestGroupBy(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}

	// Group by even/odd
	result := maps.GroupBy(input, func(key string, value int) string {
		if value%2 == 0 {
			return "even"
		}
		return "odd"
	})

	expected := map[string][]int{
		"even": {2, 4},
		"odd":  {1, 3, 5},
	}

	// Sort slices for comparison
	for k, v := range result {
		assert.ElementsMatch(t, expected[k], v)
	}
}

func TestPartition(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}

	trueMap, falseMap := maps.Partition(input, func(key string, value int) bool {
		return value%2 == 0
	})

	expectedTrue := map[string]int{"b": 2, "d": 4}
	expectedFalse := map[string]int{"a": 1, "c": 3, "e": 5}

	assert.Equal(t, expectedTrue, trueMap)
	assert.Equal(t, expectedFalse, falseMap)
}

func TestCount(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}

	// Count even values
	result := maps.Count(input, func(key string, value int) bool {
		return value%2 == 0
	})

	assert.Equal(t, 2, result)

	// Count all
	result = maps.Count(input, func(key string, value int) bool {
		return true
	})

	assert.Equal(t, 5, result)
}

func TestMin(t *testing.T) {
	input := map[string]int{"a": 3, "b": 1, "c": 5, "d": 2}

	key, value, found := maps.Min(input, func(a, b int) bool {
		return a < b
	})

	assert.True(t, found)
	assert.Equal(t, "b", key)
	assert.Equal(t, 1, value)

	// Test empty map
	_, _, found = maps.Min(map[string]int{}, func(a, b int) bool {
		return a < b
	})

	assert.False(t, found)
}

func TestMax(t *testing.T) {
	input := map[string]int{"a": 3, "b": 1, "c": 5, "d": 2}

	key, value, found := maps.Max(input, func(a, b int) bool {
		return a < b
	})

	assert.True(t, found)
	assert.Equal(t, "c", key)
	assert.Equal(t, 5, value)

	// Test empty map
	_, _, found = maps.Max(map[string]int{}, func(a, b int) bool {
		return a < b
	})

	assert.False(t, found)
}

func TestTake(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	// Take 2 elements
	result := maps.Take(input, 2)
	assert.Equal(t, 2, len(result))

	// Take more than available
	result = maps.Take(input, 10)
	assert.Equal(t, 4, len(result))
	assert.Equal(t, input, result)

	// Take 0
	result = maps.Take(input, 0)
	assert.Equal(t, 0, len(result))

	// Take negative
	result = maps.Take(input, -1)
	assert.Equal(t, 0, len(result))
}

func TestDrop(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	// Drop 2 elements
	result := maps.Drop(input, 2)
	assert.Equal(t, 2, len(result))

	// Drop more than available
	result = maps.Drop(input, 10)
	assert.Equal(t, 0, len(result))

	// Drop 0
	result = maps.Drop(input, 0)
	assert.Equal(t, 4, len(result))
	assert.Equal(t, input, result)

	// Drop negative
	result = maps.Drop(input, -1)
	assert.Equal(t, 4, len(result))
	assert.Equal(t, input, result)
}

func TestFunctionalOperationsWithDifferentTypes(t *testing.T) {
	// Test with different key/value types
	input := map[int]string{1: "one", 2: "two", 3: "three"}

	// Map
	result := maps.Map(input, func(key int, value string) int {
		return len(value)
	})
	expected := map[int]int{1: 3, 2: 3, 3: 5}
	assert.Equal(t, expected, result)

	// Filter
	filtered := maps.Filter(input, func(key int, value string) bool {
		return len(value) > 3
	})
	expectedFiltered := map[int]string{3: "three"}
	assert.Equal(t, expectedFiltered, filtered)

	// Reduce
	sum := maps.Reduce(input, 0, func(acc int, key int, value string) int {
		return acc + len(value)
	})
	assert.Equal(t, 11, sum) // 3 + 3 + 5
}

func TestFunctionalOperationsWithStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	input := map[string]Person{
		"alice":   {Name: "Alice", Age: 30},
		"bob":     {Name: "Bob", Age: 25},
		"charlie": {Name: "Charlie", Age: 35},
	}

	// Map to ages
	ages := maps.Map(input, func(key string, value Person) int {
		return value.Age
	})
	expected := map[string]int{"alice": 30, "bob": 25, "charlie": 35}
	assert.Equal(t, expected, ages)

	// Filter by age
	young := maps.Filter(input, func(key string, value Person) bool {
		return value.Age < 30
	})
	expectedYoung := map[string]Person{"bob": {Name: "Bob", Age: 25}}
	assert.Equal(t, expectedYoung, young)

	// Group by age range
	groups := maps.GroupBy(input, func(key string, value Person) string {
		if value.Age < 30 {
			return "young"
		} else if value.Age < 35 {
			return "middle"
		}
		return "senior"
	})

	assert.Equal(t, 3, len(groups))
	assert.Equal(t, 1, len(groups["young"]))
	assert.Equal(t, 1, len(groups["middle"]))
	assert.Equal(t, 1, len(groups["senior"]))
}

func TestFunctionalOperationsEdgeCases(t *testing.T) {
	// Empty map tests
	empty := map[string]int{}

	assert.Equal(t, []string{}, maps.Keys(empty))
	assert.Equal(t, []int{}, maps.Values(empty))

	// Each should not panic
	maps.Each(empty, func(key string, value int) {
		t.Fail() // Should not be called
	})

	// Map should return empty map
	result := maps.Map(empty, func(key string, value int) string {
		return "test"
	})
	assert.Equal(t, map[string]string{}, result)

	// Filter should return empty map
	result2 := maps.Filter(empty, func(key string, value int) bool {
		return true
	})
	assert.Equal(t, map[string]int{}, result2)

	// Reduce should return initial value
	result3 := maps.Reduce(empty, 42, func(acc int, key string, value int) int {
		return acc + value
	})
	assert.Equal(t, 42, result3)

	// Find should return false
	_, _, found := maps.Find(empty, func(key string, value int) bool {
		return true
	})
	assert.False(t, found)

	// Any should return false
	assert.False(t, maps.Any(empty, func(key string, value int) bool {
		return true
	}))

	// All should return true (vacuous truth)
	assert.True(t, maps.All(empty, func(key string, value int) bool {
		return false
	}))
}
