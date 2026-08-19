package uri

import (
	"net/url"
	"strings"
)

// Host returns the RFC 6454 origin of a URL (e.g. "https://server.com:8080"), lower-cased,
// dropping any userinfo, path, query, fragment, and default port. Only http and https URLs
// are recognized; anything else returns "".
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

	// Reassemble just the scheme and authority, lower-cased and canonicalized.
	return scheme + ProtocolSuffix + trimDefaultPort(scheme, strings.ToLower(parsed.Host))
}

// trimDefaultPort removes an explicit port from an authority when that port is
// the scheme's default.
func trimDefaultPort(scheme string, authority string) string {

	// RFC 6454 §6.2 omits the default port when serializing an origin, so the two
	// spellings of one origin have to produce one string. An IPv6 literal is always
	// bracketed ("[::1]:443"), which keeps the port at the tail of the authority, so
	// a suffix trim can never eat part of the host.
	switch scheme {

	case SchemeHTTP:
		return strings.TrimSuffix(authority, ":80")

	case SchemeHTTPS:
		return strings.TrimSuffix(authority, ":443")
	}

	return authority
}
