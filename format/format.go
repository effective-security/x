package format

import (
	"fmt"
)

func YesNo(val bool) string {
	if val {
		return "yes"
	}
	return "no"
}

func Enabled(val bool) string {
	if val {
		return "enabled"
	}
	return "disabled"
}

func Number[T ~int | ~int32 | ~uint | ~uint32 | ~int64 | ~uint64](val T) string {
	return fmt.Sprintf("%d", val)
}

func Float[T ~float32 | ~float64](val T) string {
	return fmt.Sprintf("%0.2f", val)
}
