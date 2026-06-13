package uri

import (
	"strings"
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

func TestHostname_EdgeCases(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Hostname(input))
	}

	// Hostnames are always lower-cased
	do("EXAMPLE.COM", "example.com")
	do("HTTPS://Example.COM:8080/Path", "example.com")

	// The port is stripped
	do("example.com:443", "example.com")
	do("https://example.com:8080", "example.com")

	// Empty input yields empty output
	do("", "")

	// Any "scheme://" protocol is stripped, not just http/https
	do("ftp://example.com", "example.com")
	do("ftp://example.com:8080/path", "example.com")
	do("ircs://chat.example.org/room", "chat.example.org")
}

func TestIsValidHostname_BoundaryLengths(t *testing.T) {

	// A label of exactly 63 characters is valid; 64 is too long (RFC 1035)
	require.True(t, IsValidHostname(strings.Repeat("a", 63)+".com"))
	require.False(t, IsValidHostname(strings.Repeat("a", 64)+".com"))

	// A domain longer than 253 characters (after punycode) is invalid
	require.False(t, IsValidHostname(strings.Repeat("a.", 130)+"com"))
}

func TestIsValidHostname_Punycode(t *testing.T) {

	// Unicode domains are converted to ASCII via punycode before validation
	require.True(t, IsValidHostname("münchen.de"))
	require.True(t, IsValidHostname("例え.jp"))
}

func TestNotValidHostname(t *testing.T) {

	// NotValidHostname is the exact inverse of IsValidHostname
	require.False(t, NotValidHostname("example.com"))
	require.False(t, NotValidHostname("localhost"))

	require.True(t, NotValidHostname(""))
	require.True(t, NotValidHostname("example.invalidtld"))
}

func TestValidateHostname(t *testing.T) {

	// Valid hostnames return no error
	require.Nil(t, ValidateHostname("example.com"))
	require.Nil(t, ValidateHostname("localhost"))
	require.Nil(t, ValidateHostname("127.0.0.1"))

	// Each distinct failure mode returns a (non-nil) validation error
	require.NotNil(t, ValidateHostname(""))                             // empty
	require.NotNil(t, ValidateHostname("com"))                          // single segment
	require.NotNil(t, ValidateHostname("example..com"))                 // empty segment
	require.NotNil(t, ValidateHostname("example.-com"))                 // invalid characters
	require.NotNil(t, ValidateHostname(strings.Repeat("a", 64)+".com")) // segment too long
	require.NotNil(t, ValidateHostname("example.invalidtld"))           // bad TLD
}

// NOTE: ValidateHostname's idna.ToASCII error branch is not covered. Inputs
// that reach it (e.g. malformed punycode) are hard to construct reliably and
// the branch is purely defensive, so it is left untested.

// FuzzHostname verifies that Hostname never panics and that its output never
// contains a path separator, a port separator, or upper-case characters.
func FuzzHostname(f *testing.F) {

	f.Add("https://example.com:8080/path?query=1")
	f.Add("EXAMPLE.COM")
	f.Add("localhost")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		result := Hostname(input)
		require.NotContains(t, result, "/", "result should not contain a slash")
		require.NotContains(t, result, ":", "result should not contain a colon")
		require.Equal(t, strings.ToLower(result), result, "result should be lower-cased")
	})
}

// FuzzValidateHostname verifies that hostname validation never panics on
// arbitrary input.
func FuzzValidateHostname(f *testing.F) {

	f.Add("example.com")
	f.Add("münchen.de")
	f.Add("")
	f.Add("...")
	f.Add(strings.Repeat("a", 300))

	f.Fuzz(func(t *testing.T, input string) {
		// IsValidHostname and NotValidHostname must always disagree.
		require.Equal(t, IsValidHostname(input), !NotValidHostname(input))
	})
}
