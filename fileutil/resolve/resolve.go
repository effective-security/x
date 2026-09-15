package resolve

import (
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/mitchellh/go-homedir"
)

// Directory returns the absolute directory name relative to baseDir.
// If the directory does not exist and create is false, a "not found" error
// is returned.
func Directory(dir string, baseDir string, create bool) (resolved string, err error) {
	if dir == "" {
		return dir, nil
	}
	if filepath.IsAbs(dir) {
		resolved = dir
	} else {
		resolved = filepath.Join(baseDir, dir)
	}
	if _, err := os.Stat(resolved); os.IsNotExist(err) {
		if create {
			if err = os.MkdirAll(resolved, 0744); err != nil {
				return "", errors.Wrapf(err, "failed to create dir: %q", resolved)
			}
		} else {
			return resolved, errors.WithStack(err)
		}
	}
	return resolved, nil
}

// File returns the absolute file name relative to baseDir.
// If the file does not exist, a "not found" error is returned.
func File(file string, baseDir string) (resolved string, err error) {
	if file == "" {
		return file, nil
	}
	if filepath.IsAbs(file) {
		resolved = file
	} else if baseDir != "" {
		resolved = filepath.Join(baseDir, file)
	} else {
		resolved = file
	}
	if _, err := os.Stat(resolved); os.IsNotExist(err) {
		return resolved, errors.WithStack(err)
	}
	return resolved, nil
}

// ExpandPath returns extrapolated path with resolved ~ or env vars.
// It replaces ${var} or $var in the string according to the values
// of the current environment variables. References to undefined
// variables are replaced by the empty string.
func ExpandPath(file string) string {
	if file == "" {
		return file
	}

	// First resolve any environment variables in the path.
	file = os.ExpandEnv(file)

	// Then resolve ~.
	file, _ = homedir.Expand(file)
	return file
}
