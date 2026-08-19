package uri

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHost_PortedFromDomainPackage(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Host(input), "input: %q", input)
	}

	do("http://localhost", "http://localhost")
	do("http://localhost/", "http://localhost")
	do("http://localhost:8080/some-path-name", "http://localhost:8080")
	do("http://localhost/some-path-name", "http://localhost")

	do("https://server.com", "https://server.com")
	do("https://server.com/", "https://server.com")
	do("https://server.com/some-path-name", "https://server.com")
	do("https://server.com/many/path/names", "https://server.com")

	// A non-default port is part of the origin and is preserved
	do("https://server.com:8443", "https://server.com:8443")

	do("not-a-host", "")
}

func TestHost_Lowercase(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Host(input), "input: %q", input)
	}

	// Scheme and host are case-insensitive, so the result is lower-cased
	do("HTTPS://CaSeInSeNsItIvE.com", "https://caseinsensitive.com")
	do("HTTP://LOCALHOST:8080/Path", "http://localhost:8080")
	do("HtTpS://Server.COM/x", "https://server.com")
}

func TestHost_StripsQueryAndFragment(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Host(input), "input: %q", input)
	}

	// Query strings and fragments are not part of the host
	do("https://server.com?query=1", "https://server.com")
	do("https://server.com/path?query=1", "https://server.com")
	do("https://server.com#fragment", "https://server.com")
	do("https://server.com/path?q=1#frag", "https://server.com")
	do("https://server.com:8080/path?q=1#frag", "https://server.com:8080")
}

func TestHost_StripsUserinfo(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Host(input), "input: %q", input)
	}

	// Userinfo (credentials) is not part of the host
	do("https://user@server.com", "https://server.com")
	do("https://user:pass@server.com", "https://server.com")
	do("https://user:pass@server.com:8080/path", "https://server.com:8080")
}

func TestHost_IPv6(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Host(input), "input: %q", input)
	}

	// IPv6 literals keep their brackets (they are part of the authority)
	do("http://[::1]:8080/path", "http://[::1]:8080")
	do("https://[2001:db8::1]/x", "https://[2001:db8::1]")
}

func TestHost_EdgeCases(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Host(input), "input: %q", input)
	}

	// Empty input yields an empty result
	do("", "")

	// Non-http(s) schemes are rejected
	do("ftp://server.com", "")
	do("mailto:ben@example.com", "")
	do("file:///etc/passwd", "")

	// A scheme with no host yields an empty result
	do("http:///path", "")
	do("https://", "")

	// A value with no scheme is not a host
	do("server.com/path", "")
	do("not-a-host", "")

	// A value that url.Parse rejects (bad %-escape, control char) is not a host
	do("https://%zz.com", "")
	do("http://server.com/\x7f", "")
}

// FuzzHost verifies that Host never panics and that its output is well-formed:
// idempotent, lower-cased, http/https only, and free of path/query/fragment.
func FuzzHost(f *testing.F) {

	f.Add("https://server.com/path")
	f.Add("http://localhost:8080")
	f.Add("HTTPS://User:Pass@Server.com:443/x?q=1#f")
	f.Add("http://[::1]:8080/path")
	f.Add("ftp://server.com")
	f.Add("not-a-host")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {

		host := Host(input)

		// Idempotent: extracting the host from a host returns the same value.
		require.Equal(t, host, Host(host), "Host is not idempotent for %q", input)

		// An empty result has no further invariants to check.
		if host == "" {
			return
		}

		// A non-empty result is lower-cased, uses a valid scheme, and contains
		// no path, query, or fragment.
		require.Equal(t, strings.ToLower(host), host, "result must be lower-cased: %q", host)
		require.True(t, strings.HasPrefix(host, ProtocolHTTP) || strings.HasPrefix(host, ProtocolHTTPS),
			"result must start with http:// or https://: %q", host)
		require.NotContains(t, host, "?", "result must not contain a query: %q", host)
		require.NotContains(t, host, "#", "result must not contain a fragment: %q", host)

		// No path: the only slashes allowed are the two in "://".
		require.Equal(t, 2, strings.Count(host, "/"), "result must not contain a path: %q", host)
	})
}

// TestHost_DefaultPort covers the RFC 6454 §6.2 rule that an origin omits the
// scheme's default port, so the two spellings of one origin agree.
func TestHost_DefaultPort(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Host(input), "input: %q", input)
	}

	// The scheme's default port is dropped
	do("https://server.com:443", "https://server.com")
	do("https://server.com:443/path?q=1", "https://server.com")
	do("http://server.com:80", "http://server.com")
	do("HTTPS://Server.COM:443/x", "https://server.com")
	do("https://[::1]:443", "https://[::1]")
	do("http://[::1]:80", "http://[::1]")

	// ...but ONLY for its own scheme; 443 is not http's default, nor 80 https'
	do("http://server.com:443", "http://server.com:443")
	do("https://server.com:80", "https://server.com:80")

	// Any other port is part of the origin
	do("https://server.com:8443", "https://server.com:8443")
	do("http://localhost:8080", "http://localhost:8080")
	do("http://[::1]:8080", "http://[::1]:8080")
}

// TestTrimDefaultPort covers the helper directly, including the schemes Host
// never passes it, so the RFC 6454 rule is pinned independent of its caller.
func TestTrimDefaultPort(t *testing.T) {

	// Each scheme drops only its OWN default port
	require.Equal(t, "server.com", trimDefaultPort(SchemeHTTP, "server.com:80"))
	require.Equal(t, "server.com:443", trimDefaultPort(SchemeHTTP, "server.com:443"))
	require.Equal(t, "server.com", trimDefaultPort(SchemeHTTPS, "server.com:443"))
	require.Equal(t, "server.com:80", trimDefaultPort(SchemeHTTPS, "server.com:80"))

	// An authority with no port is unchanged
	require.Equal(t, "server.com", trimDefaultPort(SchemeHTTPS, "server.com"))

	// An unknown scheme has no default port to drop
	require.Equal(t, "server.com:443", trimDefaultPort("ftp", "server.com:443"))
	require.Equal(t, "server.com:80", trimDefaultPort("", "server.com:80"))
}
