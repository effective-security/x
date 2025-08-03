package maps

// Keys returns all keys from the map as a slice
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values from the map as a slice
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Each applies a function to each key-value pair in the map
func Each[K comparable, V any](m map[K]V, f func(key K, value V)) {
	for k, v := range m {
		f(k, v)
	}
}

// Map transforms each value in the map using the provided function
func Map[K comparable, V any, R any](m map[K]V, f func(key K, value V) R) map[K]R {
	result := make(map[K]R, len(m))
	for k, v := range m {
		result[k] = f(k, v)
	}
	return result
}

// Filter returns a new map containing only key-value pairs that satisfy the predicate
func Filter[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// Reduce accumulates values across the map using the provided function
func Reduce[K comparable, V any, R any](m map[K]V, initial R, f func(acc R, key K, value V) R) R {
	result := initial
	for k, v := range m {
		result = f(result, k, v)
	}
	return result
}

// Find returns the first key-value pair that satisfies the predicate, or false if none found
func Find[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) (K, V, bool) {
	for k, v := range m {
		if predicate(k, v) {
			return k, v, true
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Any returns true if any key-value pair satisfies the predicate
func Any[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) bool {
	for k, v := range m {
		if predicate(k, v) {
			return true
		}
	}
	return false
}

// All returns true if all key-value pairs satisfy the predicate
func All[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) bool {
	for k, v := range m {
		if !predicate(k, v) {
			return false
		}
	}
	return true
}

// Merge combines multiple maps into one, with later maps taking precedence
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	if len(maps) == 0 {
		return make(map[K]V)
	}

	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Invert swaps keys and values (values must be comparable)
func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}

// GroupBy groups values by a key function
func GroupBy[K comparable, V any, G comparable](m map[K]V, keyFunc func(key K, value V) G) map[G][]V {
	result := make(map[G][]V)
	for k, v := range m {
		group := keyFunc(k, v)
		result[group] = append(result[group], v)
	}
	return result
}

// Partition splits the map into two based on the predicate
func Partition[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) (map[K]V, map[K]V) {
	trueMap := make(map[K]V)
	falseMap := make(map[K]V)

	for k, v := range m {
		if predicate(k, v) {
			trueMap[k] = v
		} else {
			falseMap[k] = v
		}
	}

	return trueMap, falseMap
}

// Count returns the number of key-value pairs that satisfy the predicate
func Count[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) int {
	count := 0
	for k, v := range m {
		if predicate(k, v) {
			count++
		}
	}
	return count
}

// Min returns the key-value pair with the minimum value according to the comparison function
func Min[K comparable, V any](m map[K]V, less func(a, b V) bool) (K, V, bool) {
	if len(m) == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	var minKey K
	var minValue V
	first := true

	for k, v := range m {
		if first || less(v, minValue) {
			minKey = k
			minValue = v
			first = false
		}
	}

	return minKey, minValue, true
}

// Max returns the key-value pair with the maximum value according to the comparison function
func Max[K comparable, V any](m map[K]V, less func(a, b V) bool) (K, V, bool) {
	if len(m) == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	var maxKey K
	var maxValue V
	first := true

	for k, v := range m {
		if first || less(maxValue, v) {
			maxKey = k
			maxValue = v
			first = false
		}
	}

	return maxKey, maxValue, true
}

// Take returns the first n key-value pairs from the map
func Take[K comparable, V any](m map[K]V, n int) map[K]V {
	if n <= 0 {
		return make(map[K]V)
	}

	result := make(map[K]V)
	count := 0
	for k, v := range m {
		if count >= n {
			break
		}
		result[k] = v
		count++
	}
	return result
}

// Drop returns the map with the first n key-value pairs removed
func Drop[K comparable, V any](m map[K]V, n int) map[K]V {
	if n <= 0 {
		copyMap := make(map[K]V, len(m))
		for k, v := range m {
			copyMap[k] = v
		}
		return copyMap
	}

	result := make(map[K]V)
	count := 0
	for k, v := range m {
		if count >= n {
			result[k] = v
		}
		count++
	}
	return result
}
