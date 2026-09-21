package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIsLocalURL verifies that IsLocalURL detects loopback, private, and .local hosts,
// and rejects public or unparseable URLs
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

// TestNotLocalURL verifies that NotLocalURL is the inverse of IsLocalURL
func TestNotLocalURL(t *testing.T) {

	// NotLocalURL is the exact inverse of IsLocalURL
	require.False(t, NotLocalURL("http://localhost:8080/john"))
	require.False(t, NotLocalURL("http://192.168.69.69"))

	require.True(t, NotLocalURL("https://connor.com/@john"))
	require.True(t, NotLocalURL("ftp://localhost"))
}

// TestIsLocalHostname verifies that IsLocalHostname and NotLocalHostname classify loopback,
// private, and .local hosts as local, and a public domain as not
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

// TestIsLocalHostname_EdgeCases verifies case-insensitivity, the empty string, and the
// boundaries of the 172.16/12 and 192.168 private blocks
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

// TestIsLocalHostname_ShapeTolerance verifies that IsLocalHostname finds a local host
// wrapped in a scheme, port, path, userinfo, or brackets
func TestIsLocalHostname_ShapeTolerance(t *testing.T) {

	// IsLocalHostname normalizes its input, so a "scheme://" or ":port" wrapper
	// cannot slip a local target past an SSRF check.
	require.True(t, IsLocalHostname("https://localhost/path"))
	require.True(t, IsLocalHostname("http://localhost:3000/x?y=1"))
	require.True(t, IsLocalHostname("localhost:3000"))
	require.True(t, IsLocalHostname("localhost/path"))
	require.True(t, IsLocalHostname("https://user@localhost"))
	require.True(t, IsLocalHostname("https://[::1]:8080/x"))
	require.True(t, IsLocalHostname("[::1]:8080"))
	require.True(t, IsLocalHostname("192.168.1.5:443"))
	require.True(t, IsLocalHostname("http://169.254.169.254/latest/meta-data")) // cloud metadata

	// A public host wrapped in URL syntax is still NOT local
	require.False(t, IsLocalHostname("https://example.com:8080/path"))
	require.False(t, IsLocalHostname("https://user@example.com"))
}

// TestIsLocalHostname_IPv6 verifies that IPv6 loopback, unique-local, link-local, and
// unspecified addresses are local, and a public IPv6 address is not
func TestIsLocalHostname_IPv6(t *testing.T) {

	// ::1 is the IPv6 loopback (handled as a loopback address)
	require.True(t, IsLocalHostname("::1"))

	// Unique local addresses (fc00::/7) are private/local
	require.True(t, IsLocalHostname("fc00::1"))
	require.True(t, IsLocalHostname("fd12:3456:789a:1::1"))
	require.True(t, IsLocalHostname("FD00::1")) // case-insensitive

	// Public IPv6 addresses are not local
	require.False(t, IsLocalHostname("2001:4860:4860::8888")) // public (Google DNS)

	// IPv6 link-local (fe80::/10) IS local: it is unreachable from outside this host
	// and must be blocked for SSRF safety, even though it is not a ULA.
	require.True(t, IsLocalHostname("fe80::1"))

	// IPv6 unspecified (::) is local
	require.True(t, IsLocalHostname("::"))
}

// TestIsLocalHostname_SSRFRanges verifies that link-local (including cloud metadata), unspecified,
// 127/8, and IPv4-mapped loopback and metadata addresses are local
func TestIsLocalHostname_SSRFRanges(t *testing.T) {

	// Link-local (169.254.0.0/16), including the cloud metadata endpoint
	require.True(t, IsLocalHostname("169.254.169.254")) // AWS/GCP/Azure metadata
	require.True(t, IsLocalHostname("169.254.0.1"))
	require.True(t, IsLocalHostname("169.254.255.255"))

	// Just outside 169.254.0.0/16 is public
	require.False(t, IsLocalHostname("169.253.0.1"))
	require.False(t, IsLocalHostname("169.255.0.1"))

	// Unspecified address routes to local services
	require.True(t, IsLocalHostname("0.0.0.0"))

	// Alternate loopback notations across the whole 127.0.0.0/8 block
	require.True(t, IsLocalHostname("127.0.0.2"))
	require.True(t, IsLocalHostname("127.1.2.3"))
	require.True(t, IsLocalHostname("127.255.255.254"))

	// IPv4-mapped IPv6 form of a loopback address
	require.True(t, IsLocalHostname("::ffff:127.0.0.1"))

	// IPv4-mapped IPv6 form of the metadata endpoint
	require.True(t, IsLocalHostname("::ffff:169.254.169.254"))
}

// TestIsLocalHostname_NonPublicRanges verifies that the special-use ranges IsPublicIP blocks
// (CGNAT, NAT64, documentation, benchmark, reserved) are local
func TestIsLocalHostname_NonPublicRanges(t *testing.T) {

	require.True(t, IsLocalHostname("100.64.0.1"))       // CGNAT (RFC 6598)
	require.True(t, IsLocalHostname("198.18.0.1"))       // benchmarking (RFC 2544)
	require.True(t, IsLocalHostname("192.0.2.1"))        // TEST-NET-1 documentation
	require.True(t, IsLocalHostname("240.0.0.1"))        // reserved / future use
	require.True(t, IsLocalHostname("255.255.255.255"))  // broadcast
	require.True(t, IsLocalHostname("2001:db8::1"))      // IPv6 documentation
	require.True(t, IsLocalHostname("64:ff9b::1.2.3.4")) // NAT64
	require.True(t, IsLocalHostname("2002::1"))          // 6to4

	// Genuinely public addresses remain non-local
	require.False(t, IsLocalHostname("8.8.8.8"))
	require.False(t, IsLocalHostname("2606:4700::1111"))
}

// TestIsLocalHostname_LocalStrings verifies that non-IP names for the local machine or
// network are local, and that look-alike public names are not
func TestIsLocalHostname_LocalStrings(t *testing.T) {

	// RFC 6761: "localhost" and the entire ".localhost" TLD are loopback
	require.True(t, IsLocalHostname("localhost"))
	require.True(t, IsLocalHostname("app.localhost"))
	require.True(t, IsLocalHostname("a.b.c.localhost"))

	// FQDN root form (trailing dot) is treated the same as the dotless form
	require.True(t, IsLocalHostname("localhost."))
	require.True(t, IsLocalHostname("printer.local."))

	// Common /etc/hosts loopback aliases
	require.True(t, IsLocalHostname("localhost.localdomain"))
	require.True(t, IsLocalHostname("ip6-localhost"))
	require.True(t, IsLocalHostname("ip6-loopback"))
	require.True(t, IsLocalHostname("loopback"))

	// mDNS / Bonjour and internal/Docker suffixes
	require.True(t, IsLocalHostname("printer.local"))
	require.True(t, IsLocalHostname("svc.internal"))
	require.True(t, IsLocalHostname("host.docker.internal"))

	// These must NOT be mistaken for local (substring, not suffix; or a real TLD)
	require.False(t, IsLocalHostname("example.localhostx"))
	require.False(t, IsLocalHostname("notlocalhost.com"))
	require.False(t, IsLocalHostname("myinternal.com"))
	require.False(t, IsLocalHostname("local.example.com"))
}
