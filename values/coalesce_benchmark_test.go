package values_test

import (
	"cmp"
	"testing"

	"github.com/effective-security/x/values"
)

var (
	benchmarkCoalescedString string
	benchmarkCoalescedNumber int
)

func BenchmarkStringsCoalesceVsStdlib(b *testing.B) {
	input := []string{"", "", "", "value", "unused"}
	b.Run("LegacyLoop", func(b *testing.B) {
		for b.Loop() {
			benchmarkCoalescedString = legacyStringsCoalesce(input...)
		}
	})
	b.Run("CurrentWrapper", func(b *testing.B) {
		for b.Loop() {
			benchmarkCoalescedString = values.StringsCoalesce(input...)
		}
	})
	b.Run("Stdlib", func(b *testing.B) {
		for b.Loop() {
			benchmarkCoalescedString = cmp.Or(input...)
		}
	})
}

func BenchmarkNumbersCoalesceVsStdlib(b *testing.B) {
	input := []int{0, 0, 0, 42, 99}
	b.Run("Local", func(b *testing.B) {
		for b.Loop() {
			benchmarkCoalescedNumber = values.NumbersCoalesce(input...)
		}
	})
	b.Run("Stdlib", func(b *testing.B) {
		for b.Loop() {
			benchmarkCoalescedNumber = cmp.Or(input...)
		}
	})
}

func legacyStringsCoalesce(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
