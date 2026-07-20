package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeURL(t *testing.T) {

	// Legitimate absolute http(s) URLs pass through intact (an off-site host is allowed).
	require.Equal(t, "https://remote.example/@author", SafeURL("https://remote.example/@author"))
	require.Equal(t, "http://remote.example/news", SafeURL("http://remote.example/news"))
	require.Equal(t, "https://good.com/a?b=1&c=2#frag", SafeURL("https://good.com/a?b=1&c=2#frag"))
	require.Equal(t, "https://good.com/icon.png", SafeURL("https://good.com/icon.png"))

	// Relative paths (same-origin) are safe and preserved.
	require.Equal(t, "/stream/123", SafeURL("/stream/123"))
	require.Equal(t, "", SafeURL(""))

	// Dangerous schemes are rejected outright.
	require.Equal(t, "", SafeURL("javascript:alert(1)"))
	require.Equal(t, "", SafeURL("JavaScript:alert(1)"))
	require.Equal(t, "", SafeURL("data:text/html,<script>alert(1)</script>"))
	require.Equal(t, "", SafeURL("vbscript:msgbox(1)"))
	require.Equal(t, "", SafeURL("file:///etc/passwd"))

	// Protocol-relative and leading-whitespace tricks fail closed.
	require.Equal(t, "", SafeURL("//evil.example/phish"))
	require.Equal(t, "", SafeURL("\tjavascript:alert(1)"))
}

// TestSafeURL_Neutralizes proves the returned value cannot break out of an HTML attribute of any
// quote style. Every character that could escape src="..." or onclick="f('...')" must come back
// percent-encoded, so the strings below contain no raw ", ', <, >, backtick, or space.
func TestSafeURL_Neutralizes(t *testing.T) {

	// Double-quote breakout (src="...") — the icon attack.
	require.Equal(t, "https://evil.com/x%22%20onerror=%22alert%281%29", SafeURL(`https://evil.com/x" onerror="alert(1)`))

	// Tag injection inside a URL.
	require.Equal(t, "https://evil.com/x%22%3E%3Cscript%3Ealert%281%29%3C/script%3E", SafeURL(`https://evil.com/x"><script>alert(1)</script>`))

	// Single-quote breakout (onclick="goToIntent('...')") — the id attack. url.String leaves '
	// intact, so SafeURL must encode it to %27.
	require.Equal(t, "https://evil.com/x%27),alert(1),(%27", SafeURL(`https://evil.com/x'),alert(1),('`))

	// Backtick (template-literal) is encoded by url.String.
	require.Equal(t, "https://evil.com/x%60backtick%60", SafeURL("https://evil.com/x`backtick`"))

	// A webfinger address is a safe relative value and is preserved unchanged.
	require.Equal(t, "@user@server.social", SafeURL("@user@server.social"))

	// A webfinger-shaped payload is still neutralized (no raw markup survives).
	require.NotContains(t, SafeURL(`@user@host'><script>`), "<")
	require.NotContains(t, SafeURL(`@user@host'><script>`), "'")
}
