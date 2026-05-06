package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidHostname(t *testing.T) {

	require.True(t, IsValidHostname("localhost"))
	require.True(t, IsValidHostname("friday.local"))
	require.True(t, IsValidHostname("f.local"))
	require.True(t, IsValidHostname("192.168.69.69"))
	require.True(t, IsValidHostname("something.co"))
	require.True(t, IsValidHostname("something.com"))
	require.True(t, IsValidHostname("something.org"))
	require.True(t, IsValidHostname("something.net"))
	require.True(t, IsValidHostname("something.io"))
	require.True(t, IsValidHostname("www.github.io"))

	require.False(t, IsValidHostname(""))
	require.False(t, IsValidHostname("f"))
	require.False(t, IsValidHostname("something@somewhere.com"))
	require.False(t, IsValidHostname("something.invalidtld"))
	require.False(t, IsValidHostname("something.123"))
	require.False(t, IsValidHostname("something.-com"))
	require.False(t, IsValidHostname("something.com-"))
	require.False(t, IsValidHostname("something..com"))
	require.False(t, IsValidHostname("something.com."))
	require.False(t, IsValidHostname(".something.com"))

	require.False(t, IsValidHostname("something.s"))
	require.True(t, IsValidHostname("something.so"))
	require.False(t, IsValidHostname("something.soc"))
	require.False(t, IsValidHostname("something.soci"))
	require.False(t, IsValidHostname("something.socia"))
	require.True(t, IsValidHostname("something.social"))
}

func Test_Hostname_PortedFromDomainPackage(t *testing.T) {
	require.Equal(t, "localhost", Hostname("localhost"))
	require.Equal(t, "veronica.local", Hostname("veronica.local"))
	require.Equal(t, "localhost", Hostname("https://localhost"))
	require.Equal(t, "localhost", Hostname("https://localhost/"))
	require.Equal(t, "localhost", Hostname("https://localhost/path"))
	require.Equal(t, "localhost", Hostname("https://localhost/path?and=query"))

	require.Equal(t, "localhost", Hostname("http://localhost"))
	require.Equal(t, "localhost", Hostname("http://localhost/"))
	require.Equal(t, "localhost", Hostname("http://localhost/path"))
	require.Equal(t, "localhost", Hostname("http://localhost/path?and=query"))
}
