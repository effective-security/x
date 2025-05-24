package netutil_test

import (
	"testing"

	"github.com/effective-security/x/netutil"
	"github.com/stretchr/testify/assert"
)

// TestFindFreePort_Error verifies error when no free port is found for invalid host.
func TestFindFreePort_Error(t *testing.T) {
	t.Parallel()
	p, err := netutil.FindFreePort("no-such-host", 1)
	assert.Equal(t, 0, p)
	assert.EqualError(t, err, "no free port found")
}
