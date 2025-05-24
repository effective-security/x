package configloader

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetAbsFilename_Relative verifies that GetAbsFilename resolves relative paths.
func TestGetAbsFilename_Relative(t *testing.T) {
	t.Parallel()
	file := "foo.bar"
	dir := "."
	abs, err := GetAbsFilename(file, dir)
	require.NoError(t, err)
	assert.True(t, filepath.IsAbs(abs))
	assert.Equal(t, "foo.bar", filepath.Base(abs))
}

// TestResolveConfigFile_Absolute covers the case when ResolveConfigFile is passed an absolute path.
func TestResolveConfigFile_Absolute(t *testing.T) {
	t.Parallel()
	f, err := NewFactory(nil, []string{"ignore"}, "")
	require.NoError(t, err)
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "config.yaml")
	abs, base, err := f.ResolveConfigFile(file)
	require.NoError(t, err)
	assert.Equal(t, file, abs)
	assert.Equal(t, tmpDir, base)
}

// TestResolveConfigFile_EmptyPanics verifies that ResolveConfigFile panics on empty input.
func TestResolveConfigFile_EmptyPanics(t *testing.T) {
	t.Parallel()
	f, err := NewFactory(nil, []string{"."}, "")
	require.NoError(t, err)
	assert.Panics(t, func() {
		f.ResolveConfigFile("")
	})
}
