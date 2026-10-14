package guid_test

import (
	"testing"
	"uuid"

	"github.com/effective-security/x/guid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Guid(t *testing.T) {
	g := guid.MustCreate()
	require.NotEmpty(t, g)
	assert.Equal(t, 36, len(g))

	v4 := uuid.NewV4().String()
	v7 := uuid.NewV7().String()

	assert.Equal(t, len(v4), len(g))
	assert.Equal(t, len(v7), len(g))

	v72 := uuid.NewV7().String()
	assert.Greater(t, v72, v7)
}
