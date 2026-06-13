package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidIPAddress(t *testing.T) {
	require.True(t, IsValidIPAddress("127.0.0.1"))
	require.True(t, IsValidIPAddress("::1"))
	require.False(t, IsValidIPAddress("256.256.256.256"))
	require.False(t, IsValidIPAddress("12345::6789"))
	require.False(t, IsValidIPAddress("hamburger"))
}

func TestNotValidIPAddress(t *testing.T) {

	// NotValidIPAddress is the exact inverse of IsValidIPAddress
	require.False(t, NotValidIPAddress("127.0.0.1"))
	require.False(t, NotValidIPAddress("::1"))

	require.True(t, NotValidIPAddress("256.256.256.256"))
	require.True(t, NotValidIPAddress("hamburger"))
	require.True(t, NotValidIPAddress(""))
}

func TestIsValidIP4Address(t *testing.T) {

	require.True(t, IsValidIP4Address("127.0.0.1"))
	require.True(t, IsValidIP4Address("192.168.1.1"))
	require.True(t, IsValidIP4Address("0.0.0.0"))

	// IPv6 addresses are not IPv4 addresses
	require.False(t, IsValidIP4Address("::1"))
	require.False(t, IsValidIP4Address("2001:db8::1"))

	// An IPv4-mapped IPv6 address contains colons, so it counts as IPv6, not IPv4
	require.False(t, IsValidIP4Address("::ffff:192.168.1.1"))

	// Not an IP address at all
	require.False(t, IsValidIP4Address("256.256.256.256"))
	require.False(t, IsValidIP4Address("hamburger"))
	require.False(t, IsValidIP4Address(""))
}

func TestIsValidIP6Address(t *testing.T) {

	require.True(t, IsValidIP6Address("::1"))
	require.True(t, IsValidIP6Address("2001:db8::1"))
	require.True(t, IsValidIP6Address("::ffff:192.168.1.1"))

	// IPv4 addresses are not IPv6 addresses
	require.False(t, IsValidIP6Address("127.0.0.1"))
	require.False(t, IsValidIP6Address("192.168.1.1"))

	// Not an IP address at all
	require.False(t, IsValidIP6Address("12345::6789"))
	require.False(t, IsValidIP6Address("hamburger"))
	require.False(t, IsValidIP6Address(""))
}
