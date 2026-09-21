package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestScheme verifies that Scheme returns a URL's scheme without the "://" suffix, and "" when
// the URL has no scheme or cannot be parsed
func TestScheme(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Scheme(input))
	}

	do("https://example.com", "https")
	do("http://example.com", "http")
	do("ftp://example.com", "ftp")
	do("mailto:ben@example.com", "mailto")

	// A bare hostname has no scheme
	do("example.com", "")

	// An unparseable URL returns an empty scheme
	do("https://x/%zz", "")
}

// TestIsSchemeValid verifies that only the lowercase "http" and "https" schemes are valid
func TestIsSchemeValid(t *testing.T) {

	require.True(t, IsSchemeValid("http"))
	require.True(t, IsSchemeValid("https"))

	require.False(t, IsSchemeValid(""))
	require.False(t, IsSchemeValid("ftp"))
	require.False(t, IsSchemeValid("HTTP"))     // case sensitive
	require.False(t, IsSchemeValid("https://")) // scheme does not include the suffix
}

// TestNotSchemeValid verifies that NotSchemeValid is the inverse of IsSchemeValid
func TestNotSchemeValid(t *testing.T) {

	require.False(t, NotSchemeValid("http"))
	require.False(t, NotSchemeValid("https"))

	require.True(t, NotSchemeValid(""))
	require.True(t, NotSchemeValid("ftp"))
}

// TestGuessSchemeForHostname verifies that local hostnames get SchemeHTTP and all others get SchemeHTTPS
func TestGuessSchemeForHostname(t *testing.T) {

	// Local hostnames always use plain HTTP
	require.Equal(t, SchemeHTTP, GuessSchemeForHostname("localhost"))
	require.Equal(t, SchemeHTTP, GuessSchemeForHostname("127.0.0.1"))
	require.Equal(t, SchemeHTTP, GuessSchemeForHostname("10.0.0.4"))
	require.Equal(t, SchemeHTTP, GuessSchemeForHostname("192.168.1.1"))
	require.Equal(t, SchemeHTTP, GuessSchemeForHostname("server.local"))

	// Public hostnames always use HTTPS
	require.Equal(t, SchemeHTTPS, GuessSchemeForHostname("example.com"))
	require.Equal(t, SchemeHTTPS, GuessSchemeForHostname(""))
}
