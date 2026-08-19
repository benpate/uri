package uri

import "strings"

// SafeURL accepts exactly what IsSafeRedirectURL accepts, and returns an escaped copy of
// raw that is safe to place in an HTML attribute of any quote style, or "" if raw is unsafe.
func SafeURL(raw string) string {

	// Returning "" lets the consumer omit the attribute rather than emit a dangerous one.

	// RULE: Accept exactly what IsSafeRedirectURL accepts, by asking the same helper
	// it asks. The rejection rules live in parseSafeURL and nowhere else.
	parsed, safe := parseSafeURL(raw)

	if !safe {
		return ""
	}

	// Re-serialize to percent-encode the characters that would otherwise break out of an
	// attribute. url.URL.String covers ", <, >, space, and backtick but leaves the single
	// quote intact, so encode that too -- a consumer that interpolates into src="..." or
	// onclick="f('...')" without escaping of its own is then still safe.
	return strings.ReplaceAll(parsed.String(), "'", "%27")
}
