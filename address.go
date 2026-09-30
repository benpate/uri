package uri

import "strings"

// IsSameHost returns TRUE if both values are addressed to the same host, and that host is not empty.
func IsSameHost(first string, second string) bool {

	host := AddressHost(first)

	// RULE: A value with no host matches nothing, not even another value with no host
	if host == "" {
		return false
	}

	return host == AddressHost(second)
}

// NotSameHost returns TRUE if the values are NOT addressed to the same host, including when either has no host.
func NotSameHost(first string, second string) bool {
	return !IsSameHost(first, second)
}

// AddressHost returns the lowercased host that a value is addressed to: an http(s) URL's host, the part
// of a @user@host handle (or email or acct: address) after its last "@", or empty for anything else.
func AddressHost(value string) string {

	// A URL carries its host in the authority
	if strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "http://") {
		return Hostname(value)
	}

	// An address carries its host after the last "@".  The value is classified BEFORE calling
	// NormalizeHost, which would read a bare word or "@alice.bsky.social" as a host.
	if index := strings.LastIndex(value, "@"); index > 0 {
		return NormalizeHost(value[index+1:])
	}

	return ""
}
