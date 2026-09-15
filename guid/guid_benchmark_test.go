package guid_test

import (
	"testing"
	"uuid"

	"github.com/effective-security/x/guid"
)

var benchmarkUUID string

func BenchmarkMustCreateVsStdlib(b *testing.B) {
	b.Run("LocalWrapper", func(b *testing.B) {
		for b.Loop() {
			benchmarkUUID = guid.MustCreate()
		}
	})
	b.Run("Stdlib", func(b *testing.B) {
		for b.Loop() {
			benchmarkUUID = uuid.NewV7().String()
		}
	})
}
