package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLoopback(t *testing.T) {

	// The "localhost" hostname is loopback, as is the whole RFC 6761 ".localhost"
	// TLD, the FQDN root form, and common /etc/hosts aliases
	require.True(t, IsLoopback("localhost"))
	require.True(t, IsLoopback("localhost."))
	require.True(t, IsLoopback("app.localhost"))
	require.True(t, IsLoopback("localhost.localdomain"))
	require.True(t, IsLoopback("ip6-localhost"))
	require.True(t, IsLoopback("ip6-loopback"))
	require.True(t, IsLoopback("loopback"))

	// Substring (not suffix) of ".localhost" is NOT loopback
	require.False(t, IsLoopback("example.localhostx"))

	// The entire 127.0.0.0/8 block is loopback, not just 127.0.0.1
	require.True(t, IsLoopback("127.0.0.1"))
	require.True(t, IsLoopback("127.0.0.2"))
	require.True(t, IsLoopback("127.1.2.3"))
	require.True(t, IsLoopback("127.255.255.255"))

	// IPv6 loopback and its aliases
	require.True(t, IsLoopback("::1"))
	require.True(t, IsLoopback("0:0:0:0:0:0:0:1"))  // expanded ::1
	require.True(t, IsLoopback("::ffff:127.0.0.1")) // IPv4-mapped loopback

	// Everything else is not a loopback
	require.False(t, IsLoopback(""))
	require.False(t, IsLoopback("LOCALHOST"))      // case sensitive (callers normalize first)
	require.False(t, IsLoopback("localhost:8080")) // operates on a bare host, not host:port
	require.False(t, IsLoopback("192.168.0.5"))    // private, but not loopback
	require.False(t, IsLoopback("0.0.0.0"))        // unspecified, not loopback
	require.False(t, IsLoopback("otherserver.com"))
}

func TestNotLoopback(t *testing.T) {

	// NotLoopback is the exact inverse of IsLoopback
	require.False(t, NotLoopback("localhost"))
	require.False(t, NotLoopback("127.0.0.1"))
	require.False(t, NotLoopback("::1"))

	require.True(t, NotLoopback(""))
	require.True(t, NotLoopback("192.168.0.5"))
	require.True(t, NotLoopback("otherserver.com"))
}
