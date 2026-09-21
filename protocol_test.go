package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestProtocol verifies that Protocol returns the scheme plus "://", and "" for a value that
// never had a "://" or cannot be parsed
func TestProtocol(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Protocol(input))
	}

	do("https://example.com/path", "https://")
	do("http://example.com", "http://")
	do("ftp://example.com", "ftp://")

	// A scheme with no authority still has a protocol
	do("file:///etc/passwd", "file://")

	// A value with no scheme has no protocol -- NOT a bare "://" suffix
	do("example.com", "")
	do("/just/a/path", "")
	do("", "")

	// A scheme with an opaque payload never had a "://" at all, so it has no
	// protocol -- NOT "example.com://" or "mailto://"
	do("example.com:8080/path", "")
	do("mailto:ben@example.com", "")

	// An unparseable URL returns an empty string
	do("https://x/%zz", "")
}

// TestGuessProtocolForHostname verifies that local hostnames get ProtocolHTTP and all others get ProtocolHTTPS
func TestGuessProtocolForHostname(t *testing.T) {

	// Local hostnames always use plain HTTP
	require.Equal(t, ProtocolHTTP, GuessProtocolForHostname("localhost"))
	require.Equal(t, ProtocolHTTP, GuessProtocolForHostname("192.168.1.1"))
	require.Equal(t, ProtocolHTTP, GuessProtocolForHostname("server.local"))

	// Public hostnames always use HTTPS
	require.Equal(t, ProtocolHTTPS, GuessProtocolForHostname("example.com"))
	require.Equal(t, ProtocolHTTPS, GuessProtocolForHostname(""))
}

// TestPrependProtocol verifies that PrependProtocol adds http for a local host and https
// for any other host, and returns "" when the result is not a valid http(s) URL
func TestPrependProtocol(t *testing.T) {

	do := func(input string, expected string) {
		t.Helper()
		require.Equal(t, expected, PrependProtocol(input), "PrependProtocol(%q)", input)
	}

	// URLs that already have a valid protocol are returned unchanged
	do("https://example.com", "https://example.com")
	do("http://localhost", "http://localhost")
	do("https://localhost", "https://localhost")
	do("http://example.com", "http://example.com")
	do("HTTPS://Example.com/Path", "HTTPS://Example.com/Path")

	// Public hosts get https, with or without a port or path
	do("example.com", "https://example.com")
	do("example.com:8080", "https://example.com:8080")
	do("example.com/path?q=1#frag", "https://example.com/path?q=1#frag")
	do("user@example.com", "https://user@example.com")
	do("8.8.8.8", "https://8.8.8.8")

	// Local names get http, with or without a port or path
	do("localhost", "http://localhost")
	do("localhost:8080", "http://localhost:8080")
	do("localhost:8080/path", "http://localhost:8080/path")
	do("LOCALHOST", "http://LOCALHOST")
	do("app.localhost", "http://app.localhost")
	do("friday.local", "http://friday.local")
	do("printer.local:631", "http://printer.local:631")
	do("host.docker.internal:3000", "http://host.docker.internal:3000")

	// Local and private IP addresses get http
	do("127.0.0.1", "http://127.0.0.1")
	do("127.0.0.1:8080", "http://127.0.0.1:8080")
	do("10.0.0.1", "http://10.0.0.1")
	do("192.168.1.5:8080", "http://192.168.1.5:8080")
	do("[::1]", "http://[::1]")
	do("[::1]:8080", "http://[::1]:8080")

	// Another scheme cannot be made into http(s), even when it looks like "host:port"
	do("ftp://example.com", "")
	do("mailto:me@example.com", "")
	do("javascript:alert(1)", "")
	do("tel:5551234", "")
	do("example.com:abc", "")

	// Values with no valid host produce an empty result
	do("", "")
	do("localhost:", "")
	do("//example.com", "")
	do("::1", "")
	do("not a host", "")
}

// TestIsPort verifies that isPort accepts only a non-empty run of ASCII digits
func TestIsPort(t *testing.T) {

	do := func(input string, expected bool) {
		t.Helper()
		require.Equal(t, expected, isPort(input), "isPort(%q)", input)
	}

	do("8080", true)
	do("0", true)
	do("99999999999999999999", true)

	do("", false)
	do("+80", false)
	do("-1", false)
	do("80a", false)
	do(" 80", false)
	do("٨٠", false)
}
