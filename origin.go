package uri

// IsSameOrigin returns TRUE when both values are valid http(s) URLs that share
// the same RFC 6454 origin (scheme, host, and port).
func IsSameOrigin(first string, second string) bool {

	// Comparison is origin-level, so an https URL and an http URL on the same host are
	// NOT a match, defeating a scheme-downgrade credential leak. Host canonicalizes, so
	// a URL that spells out the scheme's default port DOES match one that omits it.
	firstOrigin := Host(first)

	// RULE: An origin that cannot be established never matches.
	if firstOrigin == "" {
		return false
	}

	return firstOrigin == Host(second)
}

// NotSameOrigin returns TRUE when the two values do NOT share the same RFC 6454
// http(s) origin, including when either value is not a valid http(s) URL.
func NotSameOrigin(first string, second string) bool {
	return !IsSameOrigin(first, second)
}
