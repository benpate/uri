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

	// A port is part of the server address and is preserved
	do("https://server.com:443", "https://server.com:443")

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
