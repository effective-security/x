package netutil_test

import (
	"testing"

	"github.com/effective-security/x/netutil"
	"github.com/stretchr/testify/assert"
)

// TestIsPrivateAddress_Invalid ensures IsPrivateAddress returns error for invalid input.
func TestIsPrivateAddress_Invalid(t *testing.T) {
	t.Parallel()
	_, err := netutil.IsPrivateAddress("invalid-ip")
	assert.EqualError(t, err, "address is not valid")
}
