package ticker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNilTicker verifies nil ticker methods do not panic and return defaults.
func TestNilTicker(t *testing.T) {
	t.Parallel()
	var tkr *Ticker
	assert.Empty(t, tkr.GetStatus())
	assert.Equal(t, 0, tkr.Count())
	// Setting status on nil should be a no-op and not panic
	tkr.SetStatus("running")
	assert.Empty(t, tkr.GetStatus())
}
