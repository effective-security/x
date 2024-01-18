package values

// StringsCoalesce returns the first non-empty string value
func StringsCoalesce(str ...string) string {
	for _, s := range str {
		if len(s) > 0 {
			return s
		}
	}
	return ""
}

// NumbersCoalesce returns the first value from the supplied list that is not 0, or 0 if there are no values that are not zero
func NumbersCoalesce[T ~int | ~int32 | ~uint | ~uint32 | ~int64 | ~uint64](items ...T) T {
	for _, x := range items {
		if x != 0 {
			return x
		}
	}
	return 0
}

// Measurable interface
type Measurable[T any] interface {
	~string | ~[]string | ~[]T
}

// Coalesce returns the first non-empty value
func Coalesce[M Measurable[any]](args ...M) M {
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
