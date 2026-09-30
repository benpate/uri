package uri

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAddressHost pins which values carry a host, and which host that is
func TestAddressHost(t *testing.T) {

	test := func(value string, expected string, note string) {
		t.Helper()
		require.Equal(t, expected, AddressHost(value), "%s: %q", note, value)
	}

	// URLs
	test("https://good.example/@bob", "good.example", "https URL")
	test("http://Good.Example:8080/users/bob?x=1#main-key", "good.example", "http URL with port, query, and fragment")
	test("https://evil.example@good.example/", "good.example", "URL userinfo is not the host")
	test("https://good.example@evil.example/", "evil.example", "URL host follows the userinfo")
	test("https://[::1]:8080/x", "::1", "bracketed IPv6 URL")
	test("https://", "", "URL with an empty authority")
	test("HTTPS://good.example/", "", "scheme match is case-sensitive")
	test("ftp://good.example/", "", "non-http scheme with no address")

	// Addresses
	test("@bob@good.example", "good.example", "handle")
	test("@bob@Good.Example", "good.example", "handle is lowercased")
	test("bob@good.example", "good.example", "email address")
	test("acct:bob@good.example", "good.example", "acct: URI")
	test("bob@good.example:8080/path", "good.example", "address with port and path")
	test("a@b@good.example", "good.example", "host follows the LAST @")
	test("@bob@", "", "handle with an empty host")

	// No host
	test("#hashtag", "", "hashtag")
	test("@alice.bsky.social", "", "leading @ alone is not an address")
	test("@", "", "lone @")
	test("", "", "empty")
	test("not a key at all", "", "bare words")
}

// TestIsSameHost verifies IsSameHost and NotSameHost against matching hosts, different hosts, and hostless values
func TestIsSameHost(t *testing.T) {

	test := func(first string, second string, expected bool, note string) {
		t.Helper()
		require.Equal(t, expected, IsSameHost(first, second), "%s: %q vs %q", note, first, second)
		require.Equal(t, !expected, NotSameHost(first, second), "%s: %q vs %q", note, first, second)
	}

	const documentID = "https://good.example/@bob"

	test("https://good.example/users/bob", documentID, true, "same host, different path")
	test("https://GOOD.example/@bob", documentID, true, "host is case-insensitive")
	test("http://good.example/@bob", documentID, true, "scheme is irrelevant")
	test("https://good.example:8080/@bob", documentID, true, "port is irrelevant")
	test("@bob@good.example", documentID, true, "handle on the same host")
	test("@alice@good.example", documentID, true, "another handle on the same host")

	test("@alice@evil.example", documentID, false, "handle on another host")
	test("https://evil.example/@bob", documentID, false, "URL on another host")
	test("https://sub.good.example/@bob", documentID, false, "subdomain is another host")
	test("#hashtag", documentID, false, "hashtag has no host")
	test("@alice.bsky.social", documentID, false, "leading @ alone has no host")
	test("", documentID, false, "empty first value")
	test(documentID, "", false, "empty second value")
	test("#hashtag", "#hashtag", false, "two hostless values are not a match")
	test("", "", false, "two empty values are not a match")
}

// FuzzAddressHost verifies that AddressHost never panics, returns a lowercased host, and agrees with IsSameHost
func FuzzAddressHost(f *testing.F) {

	for _, seed := range []string{
		"https://good.example/@bob",
		"http://Good.Example:8080/users/bob?x=1#main-key",
		"https://evil.example@good.example/",
		"https://[::1]:8080/x",
		"@bob@good.example",
		"acct:bob@good.example",
		"a@b@good.example",
		"@alice.bsky.social",
		"#hashtag",
		"@",
		"",
	} {
		f.Add(seed, "https://good.example/")
	}

	f.Fuzz(func(t *testing.T, first string, second string) {

		host := AddressHost(first)

		// A host is always lowercased
		require.Equal(t, strings.ToLower(host), host, "AddressHost(%q)", first)

		// A value matches itself exactly when it has a host
		require.Equal(t, host != "", IsSameHost(first, first), "IsSameHost(%q, %q)", first, first)

		// The comparison is symmetric, and NotSameHost is its negation
		require.Equal(t, IsSameHost(first, second), IsSameHost(second, first), "symmetry of %q and %q", first, second)
		require.Equal(t, !IsSameHost(first, second), NotSameHost(first, second), "negation for %q and %q", first, second)
	})
}
