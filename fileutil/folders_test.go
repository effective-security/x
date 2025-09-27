package fileutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFolderExists(t *testing.T) {
	err := FolderExists("")
	require.Error(t, err, "expected error for empty folder")

	tmp := t.TempDir()
	f := filepath.Join(tmp, "file")
	err = os.WriteFile(f, []byte("data"), 0644)
	require.NoError(t, err, "WriteFile error")

	err = FolderExists(f)
	require.Error(t, err, "expected error for file, not a folder")

	err = FolderExists(tmp)
	require.NoError(t, err, "unexpected error for folder")
}

func TestFileExists(t *testing.T) {
	err := FileExists("")
	require.Error(t, err, "expected error for empty file")

	tmp := t.TempDir()
	err = FileExists(tmp)
	require.Error(t, err, "expected error for directory, not a file")

	f := filepath.Join(tmp, "file")
	err = os.WriteFile(f, []byte("data"), 0644)
	require.NoError(t, err, "WriteFile error")

	err = FileExists(f)
	require.NoError(t, err, "unexpected error for file")
}

func TestSubfolderNamesAndFileNames(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "sub")
	err := os.Mkdir(sub, 0755)
	require.NoError(t, err, "Mkdir error")
	file := filepath.Join(tmp, "file.txt")
	err = os.WriteFile(file, []byte("data"), 0644)
	require.NoError(t, err, "WriteFile error")

	subs, err := SubfolderNames(tmp)
	require.NoError(t, err, "SubfolderNames error")
	require.Equal(t, []string{"sub"}, subs, "unexpected subfolders")

	files, err := FileNames(tmp)
	require.NoError(t, err, "FileNames error")
	require.Equal(t, []string{"file.txt"}, files, "unexpected files")

	_, err = FileNames("not found")
	assert.EqualError(t, err, "open not found: no such file or directory")
}

func TestEnsureFolderExists(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "dir")
	err := EnsureFolderExists(dir, 0755)
	require.NoError(t, err, "EnsureFolderExists error")
	assert.DirExists(t, dir)

	file := filepath.Join(dir, "file.txt")
	err = EnsureFolderExistsForFile(file, 0755)
	require.NoError(t, err, "EnsureFolderExistsForFile error")
}
