package dns

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
