package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLoopback(t *testing.T) {

	// The three recognized loopback addresses
	require.True(t, IsLoopback("localhost"))
	require.True(t, IsLoopback("127.0.0.1"))
	require.True(t, IsLoopback("::1"))

	// Everything else is not a loopback
	require.False(t, IsLoopback(""))
	require.False(t, IsLoopback("LOCALHOST")) // case sensitive
	require.False(t, IsLoopback("localhost:8080"))
	require.False(t, IsLoopback("127.0.0.2"))
	require.False(t, IsLoopback("192.168.0.5"))
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
