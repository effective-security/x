package values

import "cmp"

// StringsCoalesce returns the first non-empty string value.
//
// Deprecated: use cmp.Or from the standard library (Go 1.22+).
func StringsCoalesce(str ...string) string {
	return cmp.Or(str...)
}

// NumbersCoalesce returns the first value from the supplied list that is not 0,
// or 0 if there are no non-zero values.
//
// Deprecated: use cmp.Or from the standard library (Go 1.22+).
func NumbersCoalesce[T ~int | ~int32 | ~uint | ~uint32 | ~int64 | ~uint64](items ...T) T {
	return cmp.Or(items...)
}

// Measurable interface
type Measurable[T any] interface {
	~string | ~[]string | ~[]T
}

// Coalesce returns the first non-empty measurable value.
// An empty argument list returns the zero value.
func Coalesce[M Measurable[any]](args ...M) M {
	var zero M
	if len(args) == 0 {
		return zero
	}
	for _, s := range args {
		if len(s) > 0 {
			return s
		}
	}
	return args[0]
}

// Select returns a if cond is true, otherwise b
func Select[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

// SelectFunc calls and returns a if cond is true; otherwise it calls and returns b.
func SelectFunc[T any](cond bool, a, b func() T) T {
	if cond {
		return a()
	}
	return b()
}
