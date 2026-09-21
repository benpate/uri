package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIsSafeRedirectURL verifies that IsSafeRedirectURL accepts relative paths and http(s) URLs,
// and rejects dangerous schemes and protocol-relative URLs
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

	// Leading whitespace is stripped before parsing, as a browser strips it,
	// so the scheme underneath is judged on its own merits and rejected.
	require.False(t, IsSafeRedirectURL("\tjavascript:alert(1)"))
	require.False(t, IsSafeRedirectURL(" javascript:alert(1)"))
}

// TestIsSafeRedirectURL_AuthorityPrefixBypass verifies that IsSafeRedirectURL and SafeURL reject
// spellings that net/url reads as a same-site path but a browser reads as an off-site authority
func TestIsSafeRedirectURL_AuthorityPrefixBypass(t *testing.T) {

	// A browser resolves every one of these to https://evil.example/x. See AGENTS.md.

	for _, target := range []string{
		`/\evil.example/x`,   // backslash is a slash for http(s)
		`\\evil.example/x`,   // both slashes backslashed
		`\/evil.example/x`,   // backslash then slash
		`/\/evil.example/x`,  // extra slashes are ignored after the authority opens
		`///evil.example/x`,  // three forward slashes, which net/url reads as an empty host
		" //evil.example/x",  // leading space, which a browser strips before it parses
		"\t//evil.example/x", // leading tab, likewise
	} {
		require.False(t, IsSafeRedirectURL(target), "target %q must be rejected", target)
		require.Equal(t, "", SafeURL(target), "target %q must not survive SafeURL", target)
	}
}

// TestIsSafeRedirectURL_BackslashInPathIsSafe verifies that IsSafeRedirectURL accepts same-site
// paths that contain a backslash or leading whitespace but do not open an authority
func TestIsSafeRedirectURL_BackslashInPathIsSafe(t *testing.T) {

	// Only the first two characters open an authority, so rejecting a later
	// backslash would break ordinary same-site paths.
	require.True(t, IsSafeRedirectURL(`/foo/\evil.example`))
	require.True(t, IsSafeRedirectURL(`\evil.example`))
	require.True(t, IsSafeRedirectURL(" /stream/123"))
}

// TestNotSafeRedirectURL verifies that NotSafeRedirectURL is the inverse of IsSafeRedirectURL
func TestNotSafeRedirectURL(t *testing.T) {

	// Unsafe targets
	require.True(t, NotSafeRedirectURL("javascript:alert(1)"))
	require.True(t, NotSafeRedirectURL("data:text/html,<script>alert(1)</script>"))
	require.True(t, NotSafeRedirectURL("//evil.com/path"))
	require.True(t, NotSafeRedirectURL("\tjavascript:alert(1)"))
	require.True(t, NotSafeRedirectURL(`/\evil.com/path`))

	// Safe targets
	require.False(t, NotSafeRedirectURL("/relative/path"))
	require.False(t, NotSafeRedirectURL("https://example.com/path"))
	require.False(t, NotSafeRedirectURL("http://example.com"))
}

// TestIsValidRedirectURL covers the stricter rule: a destination chosen by a REMOTE
// server must be an absolute http(s) URL naming a real host.
func TestIsValidRedirectURL(t *testing.T) {

	// An absolute http(s) URL is the whole point, and an off-site host is allowed.
	require.True(t, IsValidRedirectURL("https://mastodon.social/share?text=hello"))
	require.True(t, IsValidRedirectURL("http://example.com/news"))

	// Local development home servers must keep working: loopback, literal IP, and mDNS
	// names are valid hostnames even though they carry no IANA TLD.
	require.True(t, IsValidRedirectURL("http://localhost:8080/share"))
	require.True(t, IsValidRedirectURL("http://127.0.0.1:8080/share"))
	require.True(t, IsValidRedirectURL("http://friday.local/share"))

	// A relative path is not a destination here, which is the difference from
	// IsSafeRedirectURL: a remote server may not name a path on this server.
	require.False(t, IsValidRedirectURL("/signout"))
	require.False(t, IsValidRedirectURL("/@me/settings/delete"))
	require.False(t, IsValidRedirectURL("relative/path"))
	require.False(t, IsValidRedirectURL(""))

	// Everything IsSafeRedirectURL refuses stays refused.
	require.False(t, IsValidRedirectURL("javascript:alert(1)"))
	require.False(t, IsValidRedirectURL("data:text/html,<script>alert(1)</script>"))
	require.False(t, IsValidRedirectURL("file:///etc/passwd"))
	require.False(t, IsValidRedirectURL("//evil.example/x"))
	require.False(t, IsValidRedirectURL(`/\evil.example/x`))

	// A scheme with no host, or a host that cannot exist, is not a destination.
	require.False(t, IsValidRedirectURL("https://"))
	require.False(t, IsValidRedirectURL("https://phish.invalidtld/x"))
	require.False(t, IsValidRedirectURL(`https:\\evil.example\x`))
}

// TestIsValidRedirectURL_UserinfoFollowsTheHost verifies that a URL naming another site in its
// userinfo is accepted and judged by its host
func TestIsValidRedirectURL_UserinfoFollowsTheHost(t *testing.T) {

	// Off-site hosts are allowed by design, so this is a deliberate acceptance, not a bypass.
	require.True(t, IsValidRedirectURL("https://emissary.social@evil.com/login"))
	require.Equal(t, "evil.com", Hostname("https://emissary.social@evil.com/login"))
}

// TestNotValidRedirectURL verifies that NotValidRedirectURL is the inverse of IsValidRedirectURL
func TestNotValidRedirectURL(t *testing.T) {

	require.True(t, NotValidRedirectURL("/signout"))
	require.True(t, NotValidRedirectURL("javascript:alert(1)"))
	require.True(t, NotValidRedirectURL(`/\evil.example/x`))
	require.False(t, NotValidRedirectURL("https://example.com/path"))
	require.False(t, NotValidRedirectURL("http://localhost:8080/share"))
}

// TestSafeURL_BlankSentinel verifies that SafeURL returns "" both for the safe empty value
// and for "//", which IsSafeRedirectURL rejects
func TestSafeURL_BlankSentinel(t *testing.T) {

	// SafeURL's "" doubles as its rejection sentinel, so a safe value that re-serializes
	// to "" cannot be told apart from a rejection. Both are harmless as an href.
	require.True(t, IsSafeRedirectURL(""))
	require.Equal(t, "", SafeURL(""))

	// Two leading slashes open an authority, so both functions refuse "//"
	require.False(t, IsSafeRedirectURL("//"))
	require.Equal(t, "", SafeURL("//"))
}
