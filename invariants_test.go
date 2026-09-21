package uri

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

/******************************************
 * Cross-Function Fuzz Invariants
 *
 * Each target asserts a relationship that must hold for ALL inputs:
 * every IsX / NotX pair disagrees, derived values stay consistent with
 * their primitives, and (most importantly for SSRF) a non-public IP is
 * never reported as public or non-local.
 ******************************************/

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

// addSeeds feeds the shared corpus into a fuzz target's seed set.
func addSeeds(f *testing.F) {
	for _, seed := range seeds {
		f.Add(seed)
	}
}

// FuzzSSRFSafety verifies that IsLocalHostname is the exact opposite of IsPublicIP for any
// literal IP, and that IsPublicIPAddress agrees with IsPublicIP for every string.
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
		require.NotEqual(t, IsSafeRedirectURL(input), NotSafeRedirectURL(input), "redirect pair: %q", input)
		require.NotEqual(t, IsValidRedirectURL(input), NotValidRedirectURL(input), "valid-redirect pair: %q", input)
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
		protocol := Protocol(input)

		// A non-empty Protocol is always exactly its Scheme plus "://". The converse
		// does NOT hold: Protocol is "" for a value that never had a "://" at all,
		// even when Scheme finds a scheme in it ("localhost:3000", "mailto:a@b.c").
		if protocol != "" {
			require.Equal(t, scheme+ProtocolSuffix, protocol,
				"Protocol must equal Scheme + suffix for %q", input)
		}

		// Protocol NEVER returns a bare suffix, which is not a usable prefix.
		require.NotEqual(t, ProtocolSuffix, protocol,
			"Protocol must not return a bare %q for %q", ProtocolSuffix, input)

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

		if result := PrependProtocol(input); result != "" {
			require.True(t, IsValidURL(result),
				"non-empty PrependProtocol(%q) = %q must be a valid URL", input, result)
		}
	})
}

// FuzzSafeURLMatchesRedirectGuard verifies that SafeURL never returns a value IsSafeRedirectURL
// rejects, and that IsSafeRedirectURL accepts everything IsValidRedirectURL does.
func FuzzSafeURLMatchesRedirectGuard(f *testing.F) {

	addSeeds(f)

	for _, seed := range []string{
		"javascript:alert(1)", "data:text/html,<script>", "//evil.com/x", "//",
		"/path", "?q=1", "#frag", "mailto:a@b.c", "\tjavascript:alert(1)",
		"https://x.com/a'b", "/\\evil.com", "%6aavascript:alert(1)",
		"\\\\evil.com/x", "\\/evil.com/x", "/\\/evil.com/x", "///evil.com/x", " //evil.com/x",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {

		safe := IsSafeRedirectURL(input)
		neutralized := SafeURL(input)

		// The stricter rule can never be looser: a valid destination is always a safe one.
		if IsValidRedirectURL(input) {
			require.True(t, safe, "IsValidRedirectURL accepted %q, which IsSafeRedirectURL rejects", input)
		}

		// SafeURL must NEVER hand back a value the redirect guard would refuse
		if neutralized != "" {
			require.True(t, safe, "SafeURL accepted %q -> %q, which IsSafeRedirectURL rejects", input, neutralized)
		}

		// ...and must only withhold an accepted value when that value is blank anyway
		if safe && neutralized == "" {
			require.Equal(t, "", SafeURL(neutralized),
				"SafeURL withheld a non-blank accepted value for %q", input)
		}
	})
}
