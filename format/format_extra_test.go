package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSplit exercises Split for various camelcase and UTF-8 cases, including invalid UTF-8.
func TestSplit(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   string
		exp  []string
	}{
		{"Simple", "lowercase", []string{"lowercase"}},
		{"Camel", "MyClass", []string{"My", "Class"}},
		{"Acronym", "HTMLParser", []string{"HTML", "Parser"}},
		{"Digits", "Version99Test", []string{"Version", "99", "Test"}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Split(tc.in)
			assert.Equal(t, tc.exp, got)
		})
	}
}

// TestSplit_InvalidUTF8 ensures invalid UTF-8 is returned as a single entry.
func TestSplit_InvalidUTF8(t *testing.T) {
	t.Parallel()
	bad := string([]byte{0xff, 0xfe, 0xfd})
	got := Split(bad)
	assert.Equal(t, []string{bad}, got)
}
