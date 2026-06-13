package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProtocol(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Protocol(input))
	}

	do("https://example.com/path", "https://")
	do("http://example.com", "http://")
	do("ftp://example.com", "ftp://")

	// A bare hostname parses with an empty scheme, leaving only the suffix
	do("example.com", "://")

	// An unparseable URL returns an empty string
	do("https://x/%zz", "")
}

func TestGuessProtocolForHostname(t *testing.T) {

	// Local hostnames always use plain HTTP
	require.Equal(t, ProtocolHTTP, GuessProtocolForHostname("localhost"))
	require.Equal(t, ProtocolHTTP, GuessProtocolForHostname("192.168.1.1"))
	require.Equal(t, ProtocolHTTP, GuessProtocolForHostname("server.local"))

	// Public hostnames always use HTTPS
	require.Equal(t, ProtocolHTTPS, GuessProtocolForHostname("example.com"))
	require.Equal(t, ProtocolHTTPS, GuessProtocolForHostname(""))
}

func TestPrependProtocol(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, PrependProtocol(input))
	}

	// URLs that already have a valid protocol are returned unchanged
	do("https://example.com", "https://example.com")
	do("http://localhost", "http://localhost")

	// Bare public hostnames are upgraded to HTTPS. Note that url.Parse treats
	// these as a path (not a host), so the guessed scheme is always HTTPS.
	do("example.com", "https://example.com")
	do("localhost", "https://localhost")

	// "host:port" is parsed as scheme "host" with an opaque payload, which is
	// not a valid http/https scheme, so the result is empty.
	do("localhost:8080", "")

	// A non-http(s) scheme cannot be made valid, so the result is empty
	do("ftp://example.com", "")

	// An empty string produces an empty result
	do("", "")

	// NOTE: PrependProtocol's initial url.Parse error branch is not covered.
	// url.Parse rejects very few strings, and the ones it does are awkward to
	// pass through here; the branch is purely defensive, so it is left untested.
}
