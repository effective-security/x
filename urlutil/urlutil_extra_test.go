package urlutil_test

import (
	"net/http"
	"testing"

	"github.com/effective-security/x/urlutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetPublicEndpointURL_UseURLHost covers when r.URL.Host is non-empty.
func TestGetPublicEndpointURL_UseURLHost(t *testing.T) {
	t.Parallel()
	r, err := http.NewRequest(http.MethodGet, "/path", nil)
	require.NoError(t, err)

	r.URL.Scheme = "https"
	r.URL.Host = "example.com:1234"
	r.Host = "ignored.com"
	// no forwarded proto header
	u := urlutil.GetPublicEndpointURL(r, "/path").String()
	assert.Equal(t, "https://example.com:1234/path", u)
}
