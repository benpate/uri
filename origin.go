package uri

// IsSameOrigin returns TRUE when both values are valid http(s) URLs that share
// the same web origin (scheme, host, and port). Comparison is origin-level, so
// an https URL and an http URL on the same host are NOT the same origin,
// defeating a scheme-downgrade credential leak. If either value is not a valid
// http(s) URL (Host returns ""), the result is FALSE — an origin that cannot be
// established is never treated as matching.
func IsSameOrigin(first string, second string) bool {

	firstOrigin := Host(first)

	// RULE: An origin that cannot be established never matches.
	if firstOrigin == "" {
		return false
	}

	return firstOrigin == Host(second)
}

// NotSameOrigin is the inverse of IsSameOrigin. It returns TRUE when the two
// values do NOT share the same http(s) web origin (including when either value
// is not a valid http(s) URL).
func NotSameOrigin(first string, second string) bool {
	return !IsSameOrigin(first, second)
}
