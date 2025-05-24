package values

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMapAnyScanUnsupportedType covers Scan default branch for unsupported types.
func TestMapAnyScanUnsupportedType(t *testing.T) {
	t.Parallel()
	var m MapAny
	err := m.Scan(123)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported scan type:")
}

// TestIndentJSONFallback covers IndentJSON error path when input is not valid JSON.
func TestIndentJSONFallback(t *testing.T) {
	t.Parallel()
	orig := "not json"
	got := IndentJSON(orig)
	assert.Equal(t, orig, got)
}

// TestMapAnyScanInvalidJSON covers Scan path when JSON is invalid and returns unmarshal error.
func TestMapAnyScanInvalidJSON(t *testing.T) {
	t.Parallel()
	var m MapAny
	// string value that is invalid JSON
	err := m.Scan("{invalid}")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}
