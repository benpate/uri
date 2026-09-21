package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIsSameOrigin verifies IsSameOrigin and NotSameOrigin against matching origins,
// a different scheme, host, or port, and invalid values
func TestIsSameOrigin(t *testing.T) {

	testCases := []struct {
		first    string
		second   string
		expected bool
		note     string
	}{
		{"https://example.com/a", "https://example.com/b/c", true, "same host, different paths"},
		{"https://example.com", "https://example.com/", true, "path presence is irrelevant"},
		{"https://Example.COM/a", "https://example.com/b", true, "host is case-insensitive"},
		{"https://example.com:8080/a", "https://example.com:8080/b", true, "same explicit port"},

		{"https://example.com/a", "https://evil.com/b", false, "different host"},
		{"https://example.com/a", "https://sub.example.com/b", false, "subdomain is a different origin"},
		{"https://example.com/a", "http://example.com/b", false, "scheme downgrade is a different origin"},
		{"https://example.com/a", "https://example.com:8080/b", false, "different port"},

		{"", "", false, "two empty strings are not a match"},
		{"not a url", "not a url", false, "two invalid values are not a match"},
		{"https://example.com/a", "ftp://example.com/b", false, "non-http scheme is not an origin"},
		{"https://example.com/a", "", false, "empty second value"},
		{"javascript:alert(1)", "javascript:alert(1)", false, "javascript scheme is not an origin"},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.expected, IsSameOrigin(testCase.first, testCase.second), testCase.note)
		require.Equal(t, !testCase.expected, NotSameOrigin(testCase.first, testCase.second), testCase.note)
	}
}

// TestIsSameOrigin_DefaultPort covers RFC 6454: a URL that spells out the
// scheme's default port is the same origin as one that omits it.
func TestIsSameOrigin_DefaultPort(t *testing.T) {

	// Explicit default port matches the implicit one, in either position
	require.True(t, IsSameOrigin("https://example.com:443/a", "https://example.com/b"))
	require.True(t, IsSameOrigin("https://example.com/a", "https://example.com:443/b"))
	require.True(t, IsSameOrigin("http://example.com:80/a", "http://example.com/b"))
	require.True(t, IsSameOrigin("https://example.com:443/a", "https://example.com:443/b"))

	// A non-default port is still a different origin
	require.False(t, IsSameOrigin("https://example.com:8443/a", "https://example.com/b"))

	// The default port does not leak across schemes
	require.False(t, IsSameOrigin("http://example.com:443/a", "http://example.com/b"))
	require.False(t, IsSameOrigin("https://example.com:80/a", "https://example.com/b"))

	// A scheme downgrade is never the same origin, default ports or not
	require.False(t, IsSameOrigin("https://example.com:443/a", "http://example.com:80/b"))

	// NotSameOrigin stays the exact inverse
	require.False(t, NotSameOrigin("https://example.com:443/a", "https://example.com/b"))
}
