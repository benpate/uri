package uri

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIsValidHostname verifies that IsValidHostname accepts local names, IPs, and real TLDs,
// and rejects malformed labels and partial or unknown TLDs
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

// Test_Hostname_PortedFromDomainPackage verifies that Hostname returns the bare host
// of a bare hostname or an http(s) URL
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

// TestHostname_EdgeCases verifies that Hostname lower-cases its result and strips the port
// and any scheme, not just http(s)
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

// TestIsValidHostname_BoundaryLengths verifies the RFC 1035 limits of 63 characters per label
// and 253 per domain
func TestIsValidHostname_BoundaryLengths(t *testing.T) {

	// A label of exactly 63 characters is valid; 64 is too long (RFC 1035)
	require.True(t, IsValidHostname(strings.Repeat("a", 63)+".com"))
	require.False(t, IsValidHostname(strings.Repeat("a", 64)+".com"))

	// A domain longer than 253 characters (after punycode) is invalid
	require.False(t, IsValidHostname(strings.Repeat("a.", 130)+"com"))
}

// TestIsValidHostname_Punycode verifies that IsValidHostname accepts Unicode domains
func TestIsValidHostname_Punycode(t *testing.T) {

	// Unicode domains are converted to ASCII via punycode before validation
	require.True(t, IsValidHostname("münchen.de"))
	require.True(t, IsValidHostname("例え.jp"))
}

// TestNotValidHostname verifies that NotValidHostname is the inverse of IsValidHostname
func TestNotValidHostname(t *testing.T) {

	// NotValidHostname is the exact inverse of IsValidHostname
	require.False(t, NotValidHostname("example.com"))
	require.False(t, NotValidHostname("localhost"))

	require.True(t, NotValidHostname(""))
	require.True(t, NotValidHostname("example.invalidtld"))
}

// TestValidateHostname verifies that ValidateHostname returns nil for valid hostnames
// and an error for each failure mode
func TestValidateHostname(t *testing.T) {

	// Valid hostnames return no error
	require.Nil(t, ValidateHostname("example.com"))
	require.Nil(t, ValidateHostname("localhost"))
	require.Nil(t, ValidateHostname("127.0.0.1"))

	// Each distinct failure mode returns a (non-nil) validation error. The idna.ToASCII
	// error branch is not covered: it is defensive, and hard to reach reliably.
	require.NotNil(t, ValidateHostname(""))                             // empty
	require.NotNil(t, ValidateHostname("com"))                          // single segment
	require.NotNil(t, ValidateHostname("example..com"))                 // empty segment
	require.NotNil(t, ValidateHostname("example.-com"))                 // invalid characters
	require.NotNil(t, ValidateHostname(strings.Repeat("a", 64)+".com")) // segment too long
	require.NotNil(t, ValidateHostname("example.invalidtld"))           // bad TLD
}

// FuzzHostname verifies that Hostname never panics, never returns a slash or upper-case
// characters, and always strips the port from a plain "host:port".
func FuzzHostname(f *testing.F) {

	f.Add("https://example.com:8080/path?query=1")
	f.Add("EXAMPLE.COM")
	f.Add("localhost")
	f.Add("")
	f.Add("0::")
	f.Add("[::1]:8080")

	f.Fuzz(func(t *testing.T, input string) {
		result := Hostname(input)
		require.NotContains(t, result, "/", "result should not contain a slash")
		require.Equal(t, strings.ToLower(result), result, "result should be lower-cased")

		// RULE: An IPv6 literal comes back unbracketed with its colons, so only a
		// single-colon, unbracketed input is certain to be a "host:port"
		if strings.Contains(input, "[") || (strings.Count(input, ":") != 1) {
			return
		}

		require.NotContains(t, result, ":", "Hostname(%q) = %q kept the port", input, result)
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

// TestHostname_StripsQueryAndFragment asserts that a querystring or fragment never
// reaches the caller as part of the hostname.
func TestHostname_StripsQueryAndFragment(t *testing.T) {

	// A value with no path separator is the interesting shape: there is no "/" to
	// cut at, so the "?" and "#" have to be handled on their own.
	do := func(input string, expected string) {
		require.Equal(t, expected, Hostname(input), "input: %q", input)
	}

	do("example.com?q=1", "example.com")
	do("example.com#frag", "example.com")
	do("https://example.com?q=1", "example.com")
	do("https://example.com:8080?q=1#frag", "example.com")
}

// TestHostname_StripsUserinfo verifies that Hostname drops credentials from the authority
func TestHostname_StripsUserinfo(t *testing.T) {

	// Handles and referer checks are built from this result, so "evil.com@good.com"
	// must never come back as the hostname.
	do := func(input string, expected string) {
		require.Equal(t, expected, Hostname(input), "input: %q", input)
	}

	do("https://user@example.com/x", "example.com")
	do("https://user:pass@example.com:8080/x", "example.com")
	do("user@example.com", "example.com")
}

// TestHostname_IPv6 asserts that a bracketed IPv6 literal survives as a whole address.
func TestHostname_IPv6(t *testing.T) {

	// Splitting on the first colon truncates these to "[", so the brackets have to
	// be matched as a pair before any port is stripped.
	do := func(input string, expected string) {
		require.Equal(t, expected, Hostname(input), "input: %q", input)
	}

	do("[::1]:8080", "::1")
	do("http://[::1]:8080/path", "::1")
	do("https://[2001:db8::1]/x", "2001:db8::1")
}

// TestHostname_MatchesNormalizeHost verifies that Hostname and NormalizeHost return
// the same value for every shared seed.
func TestHostname_MatchesNormalizeHost(t *testing.T) {
	for _, seed := range seeds {
		require.Equal(t, NormalizeHost(seed), Hostname(seed), "input: %q", seed)
	}
}

// TestHostname_SchemeOnly verifies that a scheme with an empty authority ("https://")
// yields no hostname, and is never mistaken for one itself.
func TestHostname_SchemeOnly(t *testing.T) {

	do := func(input string, expected string) {
		t.Helper()
		require.Equal(t, expected, Hostname(input), "input: %q", input)
	}

	// A scheme with an empty authority has no host
	do("https://", "")
	do("http://", "")
	do("HTTPS://", "")
	do("ftp://", "")
	do("s3://", "")

	// ...and neither does one whose authority is empty but whose path is not
	do("https:///", "")
	do("https:///path", "")

	// RULE: "host:port" has no "//", so it must keep working
	do("localhost:8080", "localhost")
	do("example.com:443", "example.com")
}
