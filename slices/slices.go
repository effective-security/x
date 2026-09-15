// Package slices provides additional slice functions on common slice types
package slices

import (
	"cmp"
	"crypto"
	"encoding/base64"
	"regexp"
	stdslices "slices"
	"strings"

	"github.com/cockroachdb/errors"
)

// ByteSlicesEqual returns true only if the contents of the 2 slices are the same.
//
// Deprecated: use slices.Equal from the standard library.
func ByteSlicesEqual(a, b []byte) bool {
	return stdslices.Equal(a, b)
}

// StringSlicesEqual returns true only if the contents of the 2 slices are the same.
//
// Deprecated: use slices.Equal from the standard library.
func StringSlicesEqual(a, b []string) bool {
	return stdslices.Equal(a, b)
}

// ContainsString returns true if the items slice contains a value equal to item.
// Note that this can end up traversing the entire slice, and so is only really
// suitable for small slices, for larger data sets, consider using a map instead.
//
// Deprecated: use slices.Contains from the standard library.
func ContainsString(items []string, item string) bool {
	return stdslices.Contains(items, item)
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

// CloneStrings will return an independent copy of the src slice, it preserves
// the distinction between a nil value and an empty slice.
//
// Deprecated: use slices.Clone from the standard library.
func CloneStrings(src []string) []string {
	return stdslices.Clone(src)
}

// NvlString returns the first string from the supplied list that has len() > 0
// or "" if all the strings are empty.
//
// Deprecated: use cmp.Or from the standard library (Go 1.22+).
func NvlString(items ...string) string {
	return cmp.Or(items...)
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

// BoolSlicesEqual returns true only if the contents of the 2 slices are the same.
//
// Deprecated: use slices.Equal from the standard library.
func BoolSlicesEqual(a, b []bool) bool {
	return stdslices.Equal(a, b)
}

// StringUpto returns the beginning of the string up to `max`
func StringUpto(str string, maxLen int) string {
	if len(str) > maxLen {
		return str[:maxLen]
	}
	return str
}

// Int64SlicesEqual returns true only if the contents of the 2 slices are the same.
//
// Deprecated: use slices.Equal from the standard library.
func Int64SlicesEqual(a, b []int64) bool {
	return stdslices.Equal(a, b)
}

// Uint64SlicesEqual returns true only if the contents of the 2 slices are the same.
//
// Deprecated: use slices.Equal from the standard library.
func Uint64SlicesEqual(a, b []uint64) bool {
	return stdslices.Equal(a, b)
}

// Float64SlicesEqual returns true only if the contents of the 2 slices are the same.
//
// Deprecated: use slices.Equal from the standard library.
func Float64SlicesEqual(a, b []float64) bool {
	return stdslices.Equal(a, b)
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
//
// Deprecated: use slices.Contains from the standard library.
func Contains[T comparable](arr []T, val T) bool {
	return stdslices.Contains(arr, val)
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

// HashStrings returns the base64 SHA-1 of a series of string values.
// Adjacent strings are concatenated with no delimiter, so
// HashStrings("ab", "c") equals HashStrings("a", "bc").
//
// Deprecated: SHA-1 is not a safe digest. For a non-cryptographic fingerprint
// use values.XXH3HashArgs128Hex. For integrity use crypto/sha256 with explicit
// length prefixes. Do not change this function: stored hashes would break.
func HashStrings(values ...string) string {
	h := crypto.SHA1.New()
	for _, v := range values {
		_, _ = h.Write([]byte(v))
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// StringsSafeSplit splits a string into a slice of strings, removing whitespaces and empty strings.
func StringsSafeSplit(s, sep string) []string {
	list := strings.Split(s, sep)
	res := make([]string, 0, len(list))
	for _, v := range list {
		if vs := strings.TrimSpace(v); vs != "" {
			res = append(res, vs)
		}
	}
	return res
}
