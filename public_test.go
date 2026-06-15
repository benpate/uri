package uri

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsPublicIP(t *testing.T) {

	run := func(value string, expected bool) {
		ip := net.ParseIP(value)
		require.NotNil(t, ip, "unparseable IP: %s", value)
		require.Equal(t, expected, IsPublicIP(ip), "ip=%s", value)
		require.Equal(t, !expected, NotPublicIP(ip), "ip=%s", value)
	}

	// --- public addresses are allowed ---
	run("8.8.8.8", true)
	run("1.1.1.1", true)
	run("93.184.216.34", true)
	run("2606:4700:4700::1111", true)
	run("2001:4860:4860::8888", true)

	// boundaries just OUTSIDE blocked ranges are public
	run("172.15.255.255", true) // just below 172.16/12
	run("172.32.0.1", true)     // just above 172.16/12
	run("100.63.255.255", true) // just below CGNAT 100.64/10
	run("100.128.0.0", true)    // just above CGNAT 100.64/10
	run("193.168.0.1", true)    // not 192.168/16

	// --- loopback (full 127.0.0.0/8, not just 127.0.0.1) ---
	run("127.0.0.1", false)
	run("127.0.0.2", false) // the gap in uri.IsLoopback
	run("127.255.255.255", false)
	run("::1", false)

	// --- private ---
	run("10.0.0.1", false)
	run("172.16.0.1", false)
	run("172.31.255.255", false)
	run("192.168.1.1", false)
	run("fc00::1", false)
	run("fd12:3456:789a::1", false)

	// --- link-local, including the cloud metadata endpoint ---
	run("169.254.0.1", false)
	run("169.254.169.254", false) // AWS/GCP/Azure metadata
	run("fe80::1", false)

	// --- unspecified / "this network" ---
	run("0.0.0.0", false)
	run("0.1.2.3", false) // rest of 0.0.0.0/8
	run("::", false)

	// --- carrier-grade NAT ---
	run("100.64.0.1", false)
	run("100.127.255.255", false)

	// --- broadcast / reserved ---
	run("255.255.255.255", false)
	run("240.0.0.1", false)

	// --- documentation / benchmarking / protocol assignment ---
	run("192.0.0.1", false)
	run("192.0.2.1", false)
	run("198.18.0.1", false)
	run("198.51.100.1", false)
	run("203.0.113.1", false)

	// --- multicast ---
	run("224.0.0.1", false)
	run("239.255.255.255", false)
	run("ff02::1", false)

	// --- IPv6 special-use ---
	run("2001:db8::1", false)       // documentation
	run("64:ff9b::8.8.8.8", false)  // NAT64 (could embed a private v4)
	run("2002:c0a8:0101::1", false) // 6to4 embedding 192.168.1.1

	// --- IPv4-mapped IPv6 must normalize and be classified by the embedded v4 ---
	run("::ffff:127.0.0.1", false)
	run("::ffff:10.0.0.1", false)
	run("::ffff:169.254.169.254", false)
	run("::ffff:8.8.8.8", true)
}

func TestIsPublicIP_Nil(t *testing.T) {
	require.False(t, IsPublicIP(nil))
	require.True(t, NotPublicIP(nil))
}

func TestIsPublicIPAddress(t *testing.T) {

	require.True(t, IsPublicIPAddress("8.8.8.8"))
	require.True(t, IsPublicIPAddress("2606:4700:4700::1111"))

	require.False(t, IsPublicIPAddress("127.0.0.1"))
	require.False(t, IsPublicIPAddress("169.254.169.254"))
	require.False(t, IsPublicIPAddress("10.0.0.1"))

	// Unparseable strings are not public.
	require.False(t, IsPublicIPAddress(""))
	require.False(t, IsPublicIPAddress("not-an-ip"))
	require.False(t, IsPublicIPAddress("example.com"))

	// NotPublicIPAddress is the exact inverse.
	require.True(t, NotPublicIPAddress("127.0.0.1"))
	require.False(t, NotPublicIPAddress("8.8.8.8"))
}
