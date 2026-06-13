package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHost_PortedFromDomainPackage(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Host(input))
	}

	do("http://localhost", "http://localhost")
	do("http://localhost/", "http://localhost")
	do("http://localhost:8080/some-path-name", "http://localhost:8080")
	do("http://localhost/some-path-name", "http://localhost")

	do("https://server.com", "https://server.com")
	do("https://server.com/", "https://server.com")
	do("https://server.com/some-path-name", "https://server.com")
	do("https://server.com/many/path/names", "https://server.com")
	do("https://server.com:443", "https://server.com:443")

	do("HTTPS://CaSeInSeNsItIvE.com", "HTTPS://CaSeInSeNsItIvE.com")

	do("not-a-host", "")
}

func TestHost_EdgeCases(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Host(input))
	}

	// Empty input yields an empty result
	do("", "")

	// Non-http(s) schemes do not match the host regexp
	do("ftp://server.com", "")
	do("mailto:ben@example.com", "")

	// A scheme with no host after "://" does not match ([^/]+ requires one char)
	do("http:///path", "")

	// Query strings are kept (only a slash terminates the host match)
	do("https://server.com?query=1", "https://server.com?query=1")
}

// FuzzHost verifies that Host never panics and is idempotent: extracting the
// host from an already-extracted host returns the same value.
func FuzzHost(f *testing.F) {

	f.Add("https://server.com/path")
	f.Add("http://localhost:8080")
	f.Add("not-a-host")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		host := Host(input)
		require.Equal(t, host, Host(host), "Host is not idempotent for %q", input)
	})
}
