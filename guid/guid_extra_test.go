package guid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMustCreate_Panic covers the error path of MustCreate when randRead returns an error.
func TestMustCreate_Panic(t *testing.T) {
	t.Parallel()
	orig := randRead
	defer func() { randRead = orig }()
	randRead = func(b []byte) (int, error) {
		return 0, fmt.Errorf("fail")
	}
	assert.Panics(t, func() {
		MustCreate()
	})
}
