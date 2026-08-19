package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSafeRedirectURL(t *testing.T) {

	// Same-site relative paths are safe.
	require.True(t, IsSafeRedirectURL("/stream/123"))
	require.True(t, IsSafeRedirectURL("/stream/123?q=1#top"))
	require.True(t, IsSafeRedirectURL("relative/path"))
	require.True(t, IsSafeRedirectURL(""))

	// Absolute http(s) URLs are safe, even off-site.
	require.True(t, IsSafeRedirectURL("https://remote.example/@author"))
	require.True(t, IsSafeRedirectURL("http://remote.example/news"))

	// Dangerous schemes are rejected (the XSS tail on an open redirect).
	require.False(t, IsSafeRedirectURL("javascript:alert(1)"))
	require.False(t, IsSafeRedirectURL("JavaScript:alert(1)")) // scheme is case-insensitive
	require.False(t, IsSafeRedirectURL("data:text/html,<script>alert(1)</script>"))
	require.False(t, IsSafeRedirectURL("vbscript:msgbox(1)"))
	require.False(t, IsSafeRedirectURL("file:///etc/passwd"))

	// Protocol-relative URLs have an empty scheme but a host, so a browser treats
	// them as off-site absolute URLs — reject them.
	require.False(t, IsSafeRedirectURL("//evil.example/phish"))

	// Leading-whitespace tricks fail to parse and must fail closed (a browser
	// would strip the whitespace and execute the scheme).
	require.False(t, IsSafeRedirectURL("\tjavascript:alert(1)"))
}

// TestNotSafeRedirectURL confirms the Not… twin is the exact inverse of its
// Is… partner, per the package's paired-predicate convention.
func TestNotSafeRedirectURL(t *testing.T) {

	// Unsafe targets
	require.True(t, NotSafeRedirectURL("javascript:alert(1)"))
	require.True(t, NotSafeRedirectURL("data:text/html,<script>alert(1)</script>"))
	require.True(t, NotSafeRedirectURL("//evil.com/path"))
	require.True(t, NotSafeRedirectURL("\tjavascript:alert(1)"))

	// Safe targets
	require.False(t, NotSafeRedirectURL("/relative/path"))
	require.False(t, NotSafeRedirectURL("https://example.com/path"))
	require.False(t, NotSafeRedirectURL("http://example.com"))
}

// TestSafeURL_BlankSentinel documents the one place SafeURL and IsSafeRedirectURL
// appear to disagree: SafeURL's "" doubles as its rejection sentinel, so a safe
// value that re-serializes to "" cannot be told apart from a rejection. Both
// spellings are harmless as an href.
func TestSafeURL_BlankSentinel(t *testing.T) {

	for _, blank := range []string{"", "//"} {
		require.True(t, IsSafeRedirectURL(blank), "input: %q", blank)
		require.Equal(t, "", SafeURL(blank), "input: %q", blank)
	}
}
