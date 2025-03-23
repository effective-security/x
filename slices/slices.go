// Package slices provides additional slice functions on common slice types
package slices

import (
	"crypto"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
)

// ByteSlicesEqual returns true only if the contents of the 2 slices are the same
func ByteSlicesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for idx, v := range a {
		if v != b[idx] {
			return false
		}
	}
	return true
}

// StringSlicesEqual returns true only if the contents of the 2 slices are the same
func StringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for idx, v := range a {
		if v != b[idx] {
			return false
		}
	}
	return true
}

// ContainsString returns true if the items slice contains a value equal to item
// Note that this can end up traversing the entire slice, and so is only really
// suitable for small slices, for larger data sets, consider using a map instead.
func ContainsString(items []string, item string) bool {
	for _, x := range items {
		if x == item {
			return true
		}
	}
	return false
}

// StringContainsOneOf returns true if one of items slice is a substring of specified value.
func StringContainsOneOf(item string, items []string) bool {
	for _, x := range items {
		if strings.Contains(item, x) {
			return true
		}
	}
	return false
}

// StringStartsWithOneOf returns true if one of items slice is a prefix of specified value.
func StringStartsWithOneOf(value string, items []string) bool {
	for _, x := range items {
		if strings.HasPrefix(value, x) {
			return true
		}
	}
	return false
}

// ContainsStringEqualFold returns true if the items slice contains a value equal to item
// ignoring case [i.e. using EqualFold]
// Note that this can end up traversing the entire slice, and so is only really
// suitable for small slices, for larger data sets, consider using a map instead.
func ContainsStringEqualFold(items []string, item string) bool {
	for _, x := range items {
		if strings.EqualFold(x, item) {
			return true
		}
	}
	return false
}

// CloneStrings will return an independnt copy of the src slice, it preserves
// the distinction between a nil value and an empty slice.
func CloneStrings(src []string) []string {
	if src != nil {
		c := make([]string, len(src))
		copy(c, src)
		return c
	}
	return nil
}

// NvlString returns the first string from the supplied list that has len() > 0
// or "" if all the strings are empty
func NvlString(items ...string) string {
	for _, x := range items {
		if len(x) > 0 {
			return x
		}
	}
	return ""
}

// Prefixed returns a new slice of strings with each input item prefixed by the supplied prefix
// e.g. Prefixed("foo", []string{"bar","bob"}) would return []string{"foobar", "foobob"}
// the input slice is not modified.
func Prefixed(prefix string, items []string) []string {
	return MapStringSlice(items, func(in string) string {
		return prefix + in
	})
}

// Suffixed returns a new slice of strings which each input item suffixed by the supplied suffix
// e.g. Suffixed("foo", []string{"bar","bob"}) would return []string{"barfoo", "bobfoo"}
// the input slice is not modified
func Suffixed(suffix string, items []string) []string {
	return MapStringSlice(items, func(in string) string {
		return in + suffix
	})
}

// Quoted returns a new slice of strings where each input stream has been wrapped in quotes
func Quoted(items []string) []string {
	return MapStringSlice(items, func(in string) string {
		return `"` + in + `"`
	})
}

// MapStringSlice returns a new slices of strings that is the result of applies mapFn
// to each string in the input slice.
func MapStringSlice(items []string, mapFn func(in string) string) []string {
	res := make([]string, len(items))
	for idx, v := range items {
		res[idx] = mapFn(v)
	}
	return res
}

// BoolSlicesEqual returns true only if the contents of the 2 slices are the same
func BoolSlicesEqual(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for idx, v := range a {
		if v != b[idx] {
			return false
		}
	}
	return true
}

// StringUpto returns the beginning of the string up to `max`
func StringUpto(str string, maxLen int) string {
	if len(str) > maxLen {
		return str[:maxLen]
	}
	return str
}

// Int64SlicesEqual returns true only if the contents of the 2 slices are the same
func Int64SlicesEqual(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for idx, v := range a {
		if v != b[idx] {
			return false
		}
	}
	return true
}

// Uint64SlicesEqual returns true only if the contents of the 2 slices are the same
func Uint64SlicesEqual(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}
	for idx, v := range a {
		if v != b[idx] {
			return false
		}
	}
	return true
}

// Float64SlicesEqual returns true only if the contents of the 2 slices are the same
func Float64SlicesEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for idx, v := range a {
		if v != b[idx] {
			return false
		}
	}
	return true
}

// UniqueStrings removes duplicates from the given list
func UniqueStrings(dups []string) []string {
	if len(dups) < 2 {
		return dups
	}
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range dups {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Deduplicate returns a deduplicated slice.
func Deduplicate[E comparable](slice []E) []E {
	if len(slice) < 2 {
		return slice
	}

	seen := make(map[E]bool)
	deduplicated := make([]E, 0, len(slice))
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			deduplicated = append(deduplicated, v)
		}
	}
	return deduplicated
}

// Truncate returns a new slice containing the first maxLen elements of arr.
// If maxLen is greater than the length of arr, arr is returned.
func Truncate[T any](arr []T, maxLen uint) []T {
	if uint(len(arr)) <= maxLen {
		return arr
	}
	return arr[:maxLen]
}

// Contains returns true if val is in arr, and false otherwise.
func Contains[T comparable](arr []T, val T) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// StringArrayToMap converts a string array to a map.
// Each element of the string array must be in the form {key}={value}.
// - {key} is required;
// - the first '=' is the separator (also required);
// -  {value} is optional.
var reStringArrayToMap = regexp.MustCompile(`^([^=]+)=(.*)$`)

func StringArrayToMap(arr []string) (map[string]string, error) {
	m := make(map[string]string)
	for _, v := range arr {
		matches := reStringArrayToMap.FindStringSubmatch(v)
		if len(matches) != 3 {
			return nil, errors.New("invalid format for string array")
		}
		m[matches[1]] = matches[2]
	}
	return m, nil
}

// Replace replaces all occurrences of old with new in slice.
func Replace[E comparable](slice []E, old, newVal E) {
	for i, v := range slice {
		if v == old {
			slice[i] = newVal
		}
	}
}

// HashStrings returns the base64 SHA-1 of a series of string values
func HashStrings(values ...string) string {
	h := crypto.SHA1.New()
	for _, v := range values {
		_, _ = h.Write([]byte(v))
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
