package uri

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// This file collects cross-function fuzz invariants. Each target asserts a
// relationship that must hold for ALL inputs, not just hand-picked examples:
//
//   - every IsX / NotX pair must always disagree,
//   - derived values must stay consistent with their primitives,
//   - and most importantly, the SSRF-safety property: anything that resolves to
//     a non-public IP must never be reported as public/non-local.

// seeds is a shared corpus of awkward inputs worth feeding to every target.
var seeds = []string{
	"",
	" ",
	"localhost",
	"localhost:3000",
	"https://localhost/path?x=1",
	"http://127.0.0.2:8080",
	"https://user:pass@[::1]:443/x",
	"169.254.169.254",
	"http://169.254.169.254/latest/meta-data",
	"::ffff:127.0.0.1",
	"100.64.0.1",
	"2001:db8::1",
	"8.8.8.8",
	"https://example.com",
	"ftp://example.com",
	"EXAMPLE.COM:80",
	"@sarah@sky.net",
	"not a url",
	"...",
	"http://",
	"https://[::1]",
	strings.Repeat("a", 300) + ".com",
}

func addSeeds(f *testing.F) {
	for _, seed := range seeds {
		f.Add(seed)
	}
}

// FuzzSSRFSafety is the headline invariant for this package: for any literal IP,
// "local" and "public" must be exact opposites. IsLocalHostname delegates its IP
// case to NotPublicIP, so the two share a single blocklist and can never drift —
// a non-public IP (loopback, private, link-local, CGNAT, NAT64, documentation,
// benchmark, reserved, etc.) is ALWAYS classified local, and vice-versa.
func FuzzSSRFSafety(f *testing.F) {

	addSeeds(f)

	f.Fuzz(func(t *testing.T, input string) {

		host := NormalizeHost(input)

		// For any literal IP, IsLocalHostname and IsPublicIP must be complementary.
		if ip := net.ParseIP(host); ip != nil {
			require.Equal(t, IsLocalHostname(host), !IsPublicIP(ip),
				"IP %q: IsLocalHostname must be the exact opposite of IsPublicIP", host)
		}

		// IsPublicIPAddress is just IsPublicIP composed with parsing, so the two
		// must agree for every string.
		require.Equal(t, IsPublicIPAddress(host), IsPublicIP(net.ParseIP(host)),
			"IsPublicIPAddress disagrees with IsPublicIP for %q", host)
	})
}

// FuzzPredicatePairsDisagree asserts that every IsX / NotX pair is a strict
// negation of each other for all inputs.
func FuzzPredicatePairsDisagree(f *testing.F) {

	addSeeds(f)

	f.Fuzz(func(t *testing.T, input string) {

		require.NotEqual(t, IsValidHostname(input), NotValidHostname(input), "hostname pair: %q", input)
		require.NotEqual(t, IsValidIPAddress(input), NotValidIPAddress(input), "ip pair: %q", input)
		require.NotEqual(t, IsValidTLD(input), NotValidTLD(input), "tld pair: %q", input)
		require.NotEqual(t, IsValidURL(input), NotValidURL(input), "url pair: %q", input)
		require.NotEqual(t, IsSchemeValid(input), NotSchemeValid(input), "scheme pair: %q", input)
		require.NotEqual(t, IsLoopback(input), NotLoopback(input), "loopback pair: %q", input)
		require.NotEqual(t, IsLocalHostname(input), NotLocalHostname(input), "local-hostname pair: %q", input)
		require.NotEqual(t, IsLocalURL(input), NotLocalURL(input), "local-url pair: %q", input)
		require.NotEqual(t, IsPublicIPAddress(input), NotPublicIPAddress(input), "public-ip pair: %q", input)
	})
}

// FuzzValidateMatchesPredicate asserts that each ValidateX returns nil exactly
// when its IsX predicate is true.
func FuzzValidateMatchesPredicate(f *testing.F) {

	addSeeds(f)

	f.Fuzz(func(t *testing.T, input string) {

		require.Equal(t, ValidateHostname(input) == nil, IsValidHostname(input), "hostname: %q", input)
		require.Equal(t, ValidateURL(input) == nil, IsValidURL(input), "url: %q", input)
		require.Equal(t, ValidateTLD(input) == nil, IsValidTLD(input), "tld: %q", input)
	})
}

// FuzzSchemeProtocolConsistency asserts the structural relationship between a
// URL's scheme and its protocol, and that a valid scheme round-trips.
func FuzzSchemeProtocolConsistency(f *testing.F) {

	addSeeds(f)

	f.Fuzz(func(t *testing.T, input string) {

		scheme := Scheme(input)

		// When a scheme is present, Protocol is exactly that scheme plus "://".
		// (The empty-scheme case is intentionally not asserted here: Scheme and
		// Protocol diverge on url.Parse errors — see the package notes on the
		// Scheme/Protocol empty-input inconsistency.)
		if scheme != "" {
			require.Equal(t, scheme+ProtocolSuffix, Protocol(input),
				"Protocol must equal Scheme + suffix for %q", input)
		}

		// IsSchemeValid only accepts the two known schemes, and validity implies
		// the scheme is exactly one of them.
		if IsSchemeValid(scheme) {
			require.True(t, scheme == SchemeHTTP || scheme == SchemeHTTPS,
				"valid scheme %q must be http or https", scheme)
		}
	})
}

// FuzzIPAddressClassUnion asserts that an address is valid if and only if it is
// classified as either IPv4 or IPv6 (the two classes partition validity).
func FuzzIPAddressClassUnion(f *testing.F) {

	addSeeds(f)

	f.Fuzz(func(t *testing.T, input string) {
		require.Equal(t, IsValidIPAddress(input), IsValidIP4Address(input) || IsValidIP6Address(input),
			"IsValidIPAddress must equal IP4||IP6 for %q", input)
	})
}

// FuzzPrependProtocolWellFormed asserts that PrependProtocol either fails
// (returns "") or returns a string that is itself a valid URL.
func FuzzPrependProtocolWellFormed(f *testing.F) {

	addSeeds(f)

	f.Fuzz(func(t *testing.T, input string) {

		result := PrependProtocol(input)

		if result != "" {
			require.True(t, IsValidURL(result),
				"non-empty PrependProtocol(%q) = %q must be a valid URL", input, result)
		}
	})
}
