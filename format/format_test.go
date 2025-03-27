package format_test

import (
	"testing"

	"github.com/effective-security/x/format"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYesNo(t *testing.T) {
	assert.Equal(t, "yes", format.YesNo(true))
	assert.Equal(t, "no", format.YesNo(false))
}

func TestEnabled(t *testing.T) {
	assert.Equal(t, "enabled", format.Enabled(true))
	assert.Equal(t, "disabled", format.Enabled(false))
}

func TestNumber(t *testing.T) {
	require.Equal(t, "42", format.Number(42))
	require.Equal(t, "123456789", format.Number(int64(123456789)))
	require.Equal(t, "0", format.Number(uint(0)))
}

func TestFloat(t *testing.T) {
	require.Equal(t, "3.14", format.Float(3.14159))
	require.Equal(t, "0.00", format.Float(0.0))
	require.Equal(t, "-1.23", format.Float(-1.234))
}
