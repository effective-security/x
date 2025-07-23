package values

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncMapStringInt(t *testing.T) {
	m := SyncMap[string, int]{}

	// Test Store and Load
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	val, ok := m.Load("a")
	require.True(t, ok)
	assert.Equal(t, 1, val)

	val, ok = m.Load("b")
	require.True(t, ok)
	assert.Equal(t, 2, val)

	val, ok = m.Load("nonexistent")
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	// Test LoadOrStore
	val, loaded := m.LoadOrStore("a", 10)
	require.True(t, loaded)
	assert.Equal(t, 1, val)

	val, loaded = m.LoadOrStore("d", 4)
	assert.False(t, loaded)
	assert.Equal(t, 4, val)

	// Test LoadAndDelete
	val, loaded = m.LoadAndDelete("b")
	require.True(t, loaded)
	assert.Equal(t, 2, val)

	val, loaded = m.LoadAndDelete("b")
	assert.False(t, loaded)
	assert.Equal(t, 0, val)

	// Test Delete
	m.Delete("c")
	val, ok = m.Load("c")
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	// Test Range
	expected := map[string]int{"a": 1, "d": 4}
	actual := make(map[string]int)
	m.Range(func(key string, value int) bool {
		actual[key] = value
		return true
	})

	assert.Equal(t, expected, actual)
}

func TestSyncMapIntString(t *testing.T) {
	m := SyncMap[int, string]{}

	m.Store(1, "one")
	m.Store(2, "two")
	m.Store(3, "three")

	val, ok := m.Load(1)
	require.True(t, ok)
	assert.Equal(t, "one", val)

	// Test LoadOrStore with existing key
	val, loaded := m.LoadOrStore(1, "ONE")
	require.True(t, loaded)
	assert.Equal(t, "one", val)

	// Test LoadOrStore with new key
	val, loaded = m.LoadOrStore(4, "four")
	assert.False(t, loaded)
	assert.Equal(t, "four", val)

	// Test LoadAndDelete
	val, loaded = m.LoadAndDelete(2)
	require.True(t, loaded)
	assert.Equal(t, "two", val)

	m.Delete(3)

	expected := map[int]string{1: "one", 4: "four"}
	actual := make(map[int]string)
	m.Range(func(key int, value string) bool {
		actual[key] = value
		return true
	})

	assert.Equal(t, expected, actual)
}

func TestSyncMapFloat64Bool(t *testing.T) {
	m := SyncMap[float64, bool]{}

	m.Store(1.5, true)
	m.Store(2.7, false)
	m.Store(3.14, true)

	val, ok := m.Load(1.5)
	require.True(t, ok)
	assert.True(t, val)

	val, ok = m.Load(2.7)
	require.True(t, ok)
	assert.False(t, val)

	// Test LoadOrStore
	val, loaded := m.LoadOrStore(1.5, false)
	require.True(t, loaded)
	assert.True(t, val)

	val, loaded = m.LoadOrStore(4.2, false)
	assert.False(t, loaded)
	assert.False(t, val)

	// Test LoadAndDelete
	val, loaded = m.LoadAndDelete(2.7)
	require.True(t, loaded)
	assert.False(t, val)

	m.Delete(3.14)

	expected := map[float64]bool{1.5: true, 4.2: false}
	actual := make(map[float64]bool)
	m.Range(func(key float64, value bool) bool {
		actual[key] = value
		return true
	})

	assert.Equal(t, expected, actual)
}

func TestSyncMapStructTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	type Address struct {
		Street  string
		City    string
		Country string
	}

	m := SyncMap[Person, Address]{}

	person1 := Person{Name: "Alice", Age: 30}
	person2 := Person{Name: "Bob", Age: 25}

	addr1 := Address{Street: "123 Main St", City: "New York", Country: "USA"}
	addr2 := Address{Street: "456 Oak Ave", City: "London", Country: "UK"}

	m.Store(person1, addr1)
	m.Store(person2, addr2)

	val, ok := m.Load(person1)
	require.True(t, ok)
	assert.Equal(t, addr1, val)

	// Test LoadOrStore
	val, loaded := m.LoadOrStore(person1, Address{})
	require.True(t, loaded)
	assert.Equal(t, addr1, val)

	person3 := Person{Name: "Charlie", Age: 35}
	addr3 := Address{Street: "789 Pine Rd", City: "Paris", Country: "France"}

	val, loaded = m.LoadOrStore(person3, addr3)
	assert.False(t, loaded)
	assert.Equal(t, addr3, val)

	// Test LoadAndDelete
	val, loaded = m.LoadAndDelete(person2)
	require.True(t, loaded)
	assert.Equal(t, addr2, val)

	m.Delete(person1)

	expected := map[Person]Address{person3: addr3}
	actual := make(map[Person]Address)
	m.Range(func(key Person, value Address) bool {
		actual[key] = value
		return true
	})

	assert.Equal(t, expected, actual)
}

func TestSyncMapSliceTypes(t *testing.T) {
	m := SyncMap[string, []int]{}

	m.Store("evens", []int{2, 4, 6, 8})
	m.Store("odds", []int{1, 3, 5, 7})
	m.Store("primes", []int{2, 3, 5, 7})

	val, ok := m.Load("evens")
	require.True(t, ok)
	assert.Equal(t, []int{2, 4, 6, 8}, val)

	// Test LoadOrStore
	val, loaded := m.LoadOrStore("evens", []int{10, 12})
	require.True(t, loaded)
	assert.Equal(t, []int{2, 4, 6, 8}, val)

	val, loaded = m.LoadOrStore("fibonacci", []int{1, 1, 2, 3, 5})
	assert.False(t, loaded)
	assert.Equal(t, []int{1, 1, 2, 3, 5}, val)

	// Test LoadAndDelete
	val, loaded = m.LoadAndDelete("odds")
	require.True(t, loaded)
	assert.Equal(t, []int{1, 3, 5, 7}, val)

	m.Delete("primes")

	expected := map[string][]int{
		"evens":     {2, 4, 6, 8},
		"fibonacci": {1, 1, 2, 3, 5},
	}
	actual := make(map[string][]int)
	m.Range(func(key string, value []int) bool {
		actual[key] = value
		return true
	})

	assert.Equal(t, expected, actual)
}

func TestSyncMapEmptyMap(t *testing.T) {
	m := SyncMap[string, int]{}

	// Test Load on empty map
	val, ok := m.Load("any")
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	// Test LoadAndDelete on empty map
	val, loaded := m.LoadAndDelete("any")
	assert.False(t, loaded)
	assert.Equal(t, 0, val)

	// Test Range on empty map
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return true
	})

	assert.Equal(t, 0, count)
}

func TestSyncMapRangeEarlyReturn(t *testing.T) {
	m := SyncMap[int, string]{}

	for i := 1; i <= 10; i++ {
		m.Store(i, fmt.Sprintf("value-%d", i))
	}

	// Test Range with early return
	count := 0
	m.Range(func(key int, value string) bool {
		count++
		return count < 5 // Stop after 4 iterations
	})

	// The Range function processes the current item before checking the return value
	// So it will process 5 items and then stop
	assert.Equal(t, 5, count)
}

func TestSyncMapConcurrentAccess(t *testing.T) {
	m := SyncMap[int, string]{}
	done := make(chan bool)

	// Start goroutines to concurrently access the map
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				key := id*100 + j
				m.Store(key, fmt.Sprintf("value-%d", key))
				m.Load(key)
				if j%10 == 0 {
					m.Delete(key)
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify the map still works correctly
	count := 0
	m.Range(func(key int, value string) bool {
		count++
		return true
	})

	// Should have some remaining items (not all were deleted)
	assert.Greater(t, count, 0, "Map should not be empty after concurrent access")
}

func TestSyncMapPointerTypes(t *testing.T) {
	m := SyncMap[string, *int]{}

	val1 := 42
	val2 := 100

	m.Store("ptr1", &val1)
	m.Store("ptr2", &val2)

	ptr, ok := m.Load("ptr1")
	require.True(t, ok)
	assert.Equal(t, &val1, ptr)
	assert.Equal(t, 42, *ptr)

	// Test LoadOrStore
	ptr, loaded := m.LoadOrStore("ptr1", &val2)
	require.True(t, loaded)
	assert.Equal(t, &val1, ptr)

	val3 := 200
	ptr, loaded = m.LoadOrStore("ptr3", &val3)
	assert.False(t, loaded)
	assert.Equal(t, &val3, ptr)

	// Test LoadAndDelete
	ptr, loaded = m.LoadAndDelete("ptr2")
	require.True(t, loaded)
	assert.Equal(t, &val2, ptr)

	m.Delete("ptr1")

	expected := map[string]*int{"ptr3": &val3}
	actual := make(map[string]*int)
	m.Range(func(key string, value *int) bool {
		actual[key] = value
		return true
	})

	assert.Equal(t, expected, actual)
}

func TestSyncMapInterfaceTypes(t *testing.T) {
	m := SyncMap[string, interface{}]{}

	m.Store("string", "hello")
	m.Store("int", 42)
	m.Store("bool", true)
	m.Store("slice", []int{1, 2, 3})

	val, ok := m.Load("string")
	require.True(t, ok)
	assert.Equal(t, "hello", val)

	val, ok = m.Load("int")
	require.True(t, ok)
	assert.Equal(t, 42, val)

	val, ok = m.Load("bool")
	require.True(t, ok)
	assert.Equal(t, true, val)

	val, ok = m.Load("slice")
	require.True(t, ok)
	assert.Equal(t, []int{1, 2, 3}, val)

	// Test LoadOrStore
	val, loaded := m.LoadOrStore("string", "world")
	require.True(t, loaded)
	assert.Equal(t, "hello", val)

	val, loaded = m.LoadOrStore("float", 3.14)
	assert.False(t, loaded)
	assert.Equal(t, 3.14, val)

	// Test LoadAndDelete
	val, loaded = m.LoadAndDelete("int")
	require.True(t, loaded)
	assert.Equal(t, 42, val)

	m.Delete("bool")

	expected := map[string]interface{}{
		"string": "hello",
		"slice":  []int{1, 2, 3},
		"float":  3.14,
	}
	actual := make(map[string]interface{})
	m.Range(func(key string, value interface{}) bool {
		actual[key] = value
		return true
	})

	assert.Equal(t, expected, actual)
}

func TestSyncMapZeroValueHandling(t *testing.T) {
	m := SyncMap[string, int]{}

	// Test that zero values are handled correctly
	val, ok := m.Load("nonexistent")
	assert.False(t, ok)
	assert.Equal(t, 0, val) // zero value for int

	val, loaded := m.LoadAndDelete("nonexistent")
	assert.False(t, loaded)
	assert.Equal(t, 0, val)

	// Test with bool type
	mBool := SyncMap[string, bool]{}
	valBool, ok := mBool.Load("nonexistent")
	assert.False(t, ok)
	assert.False(t, valBool) // zero value for bool

	// Test with string type
	mStr := SyncMap[string, string]{}
	valStr, ok := mStr.Load("nonexistent")
	assert.False(t, ok)
	assert.Equal(t, "", valStr) // zero value for string
}
