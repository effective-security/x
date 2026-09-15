package resolve_test

import (
	"os"
	"testing"

	"github.com/effective-security/x/fileutil/resolve"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestFile_EmptyBaseDirUsesPath(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.WriteFile("x.txt", []byte("ok"), 0o600))
	got, err := resolve.File("x.txt", "")
	require.NoError(t, err)
	assert.Equal(t, "x.txt", got)
}
