package netutil

import (
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsAddrInUse_True ensures IsAddrInUse detects EADDRINUSE errors.
func TestIsAddrInUse_True(t *testing.T) {
	t.Parallel()
	err := &net.OpError{Err: &os.SyscallError{Err: syscall.EADDRINUSE}}
	assert.True(t, IsAddrInUse(err))
}

// TestNewNamedAddress_Error covers newNamedAddress error path for unresolved host.
func TestNewNamedAddress_Error(t *testing.T) {
	t.Parallel()
	_, err := newNamedAddress("tcp", "no-such-host", 1234)
	assert.Error(t, err)
}

// TestIsAddrInUse_FalseCases checks that non EADDRINUSE errors return false.
func TestIsAddrInUse_FalseCases(t *testing.T) {
	t.Parallel()
	assert.False(t, IsAddrInUse(nil))
	assert.False(t, IsAddrInUse(assert.AnError))
}
