package resolve_test

import (
	"testing"

	"github.com/effective-security/x/fileutil/resolve"
	"github.com/stretchr/testify/assert"
)

// TestDirectoryEmpty verifies that an empty dir returns empty without error.
func TestDirectoryEmpty(t *testing.T) {
	t.Parallel()
	d, err := resolve.Directory("", "/base", false)
	assert.NoError(t, err)
	assert.Equal(t, "", d)

	d, err = resolve.Directory("", "/base", true)
	assert.NoError(t, err)
	assert.Equal(t, "", d)
}
