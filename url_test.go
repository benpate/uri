package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidURL(t *testing.T) {

	// Woot woot.
	require.True(t, IsValidURL("http://example.com"))

	// Invalid hostnames
	require.False(t, IsValidURL("http://example.socia"))

	// Invalid protocols
	require.False(t, IsValidURL("ftp://example.com"))
	require.False(t, IsValidURL("example.com"))
}
