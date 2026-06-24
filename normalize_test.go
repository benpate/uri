package uri

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeHost(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, NormalizeHost(input), "input: %q", input)
	}

	// Bare hostnames pass through (lower-cased)
	do("localhost", "localhost")
	do("example.com", "example.com")
	do("EXAMPLE.com", "example.com")
	do("veronica.local", "veronica.local")

	// "host:port" has its port stripped, even without a scheme
	do("localhost:3000", "localhost")
	do("example.com:443", "example.com")
	do("EXAMPLE.com:80", "example.com")
	do("192.168.1.5:443", "192.168.1.5")

	// Full URLs of any scheme yield just the host
	do("https://localhost/path", "localhost")
	do("http://localhost:3000/x?y=1", "localhost")
	do("ftp://localhost/x", "localhost")
	do("HTTPS://LocalHost/PATH", "localhost")

	// No-scheme values with a path still resolve the host
	do("localhost/path", "localhost")

	// Userinfo is stripped, with or without a scheme
	do("https://user@localhost", "localhost")
	do("user:pass@localhost:8080", "localhost")
	do("@sarah@sky.net", "sky.net")

	// IPv6 literals are returned unbracketed (and MAY contain colons)
	do("https://[::1]:8080/x", "::1")
	do("[::1]:8080", "::1")
	do("::1", "::1")
	do("[2001:db8::1]", "2001:db8::1")

	// Trailing dots and empty input
	do("internal.local.", "internal.local.")
	do("https://internal.local./x", "internal.local.")
	do("", "")
}

// FuzzNormalizeHost verifies that NormalizeHost never panics and always
// lower-cases its output. (Unlike Hostname, the output MAY contain colons,
// so that invariant is intentionally NOT asserted here.)
func FuzzNormalizeHost(f *testing.F) {

	f.Add("https://example.com:8080/path?query=1")
	f.Add("user:pass@[::1]:443")
	f.Add("EXAMPLE.COM")
	f.Add("localhost")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		result := NormalizeHost(input)
		require.Equal(t, strings.ToLower(result), result, "result should be lower-cased")
		require.NotContains(t, result, "/", "result should not contain a slash")
	})
}
