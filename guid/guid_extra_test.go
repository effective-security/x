package guid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMustCreate_Panic covers the error path of MustCreate when randRead returns an error.
func TestMustCreate_Panic(t *testing.T) {
	t.Parallel()
	assert.NotPanics(t, func() {
		MustCreate()
	})
}
