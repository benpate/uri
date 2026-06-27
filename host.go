package uri

import (
	"net/url"
	"strings"
)

// Host returns the scheme and authority of a URL (e.g. "https://server.com:8080"),
// dropping any userinfo, path, query, and fragment. The scheme and host are
// lower-cased. Only http and https URLs are recognized; anything else (an
// unparseable value, a missing scheme or host, or another scheme) returns "".
func Host(value string) string {

	// Parse the URL. An unparseable value is not a host.
	parsed, err := url.Parse(value)

	if err != nil {
		return ""
	}

	// RULE: Only http and https URLs have a host we recognize.
	scheme := strings.ToLower(parsed.Scheme)

	if (scheme != SchemeHTTP) && (scheme != SchemeHTTPS) {
		return ""
	}

	// RULE: A URL with no host (e.g. "http:///path") is not a host.
	// parsed.Host includes the port and a bracketed IPv6 literal, but not userinfo.
	if parsed.Host == "" {
		return ""
	}

	// Reassemble just the scheme and authority, lower-cased.
	return scheme + ProtocolSuffix + strings.ToLower(parsed.Host)
}
