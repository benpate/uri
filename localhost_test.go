package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLocalURL(t *testing.T) {

	require.True(t, IsLocalURL("http://127.0.0.1/john"))
	require.True(t, IsLocalURL("http://10.0.0.4/@john"))
	require.True(t, IsLocalURL("http://localhost:8080/john"))
	require.True(t, IsLocalURL("http://192.168.69.69"))
	require.True(t, IsLocalURL("http://172.16.69.69"))
	require.True(t, IsLocalURL("https://server.local"))

	require.False(t, IsLocalURL("http://connor.com"))
	require.False(t, IsLocalURL("https://connor.com/@john"))
	require.False(t, IsLocalURL("https://connor.com:1234/@john"))

	// A URL that fails ParseURL (bad scheme) is not considered local
	require.False(t, IsLocalURL("ftp://localhost"))
	require.False(t, IsLocalURL("not a url"))
}

func TestNotLocalURL(t *testing.T) {

	// NotLocalURL is the exact inverse of IsLocalURL
	require.False(t, NotLocalURL("http://localhost:8080/john"))
	require.False(t, NotLocalURL("http://192.168.69.69"))

	require.True(t, NotLocalURL("https://connor.com/@john"))
	require.True(t, NotLocalURL("ftp://localhost"))
}

func TestIsLocalHostname(t *testing.T) {

	require.True(t, IsLocalHostname("localhost"))
	require.False(t, NotLocalHostname("localhost"))

	require.True(t, IsLocalHostname("127.0.0.1"))
	require.True(t, IsLocalHostname("10.0.0.4"))
	require.True(t, IsLocalHostname("192.168.69.69"))
	require.True(t, IsLocalHostname("172.16.69.69"))
	require.True(t, IsLocalHostname("server.local"))

	require.False(t, IsLocalHostname("connor.com"))
	require.True(t, NotLocalHostname("connor.com"))
}

func TestIsLocalHostname_EdgeCases(t *testing.T) {

	// The check is case-insensitive
	require.True(t, IsLocalHostname("LOCALHOST"))
	require.True(t, IsLocalHostname("Server.LOCAL"))

	// Any ".local" suffix is treated as local
	require.True(t, IsLocalHostname("anything.local"))

	// Empty string is not local
	require.False(t, IsLocalHostname(""))

	// The full 172.16.0.0/12 block (172.16.0.0 - 172.31.255.255) is local
	require.True(t, IsLocalHostname("172.16.0.1"))
	require.True(t, IsLocalHostname("172.17.0.1"))
	require.True(t, IsLocalHostname("172.31.255.255"))

	// Addresses just outside the 172.16/12 block are public
	require.False(t, IsLocalHostname("172.15.0.1"))
	require.False(t, IsLocalHostname("172.32.0.1"))

	// "172.160.x" must NOT be mistaken for the "172.16" block
	require.False(t, IsLocalHostname("172.160.0.1"))

	// "192.168" must be the prefix; "193.168" is public
	require.True(t, IsLocalHostname("192.168.0.1"))
	require.False(t, IsLocalHostname("193.168.0.1"))

	// A public IP is not local
	require.False(t, IsLocalHostname("8.8.8.8"))
}

func TestIsLocalHostname_IPv6(t *testing.T) {

	// ::1 is the IPv6 loopback (handled as a loopback address)
	require.True(t, IsLocalHostname("::1"))

	// Unique local addresses (fc00::/7) are private/local
	require.True(t, IsLocalHostname("fc00::1"))
	require.True(t, IsLocalHostname("fd12:3456:789a:1::1"))
	require.True(t, IsLocalHostname("FD00::1")) // case-insensitive

	// Public and link-local IPv6 addresses are not "private"
	require.False(t, IsLocalHostname("2001:4860:4860::8888")) // public (Google DNS)
	require.False(t, IsLocalHostname("fe80::1"))              // link-local, not ULA
}
