package uri

import (
	"net/url"
	"strings"
)

/******************************************
 * Protocols
 *
 * A Protocol is the "protocol" portion of a
 * URL, which is the connection scheme plus
 * the "://" suffix.
 ******************************************/

// Protocol returns the protocol portion of a URL.
// If the URL is not valid, or has no scheme, then an empty string is returned.
func Protocol(uri string) string {

	parsedURL, err := url.Parse(uri)

	if err != nil {
		return ""
	}

	// RULE: A value with no scheme has no protocol. Returning a bare "://" would
	// be a prefix that no caller can concatenate into a valid URL.
	if parsedURL.Scheme == "" {
		return ""
	}

	// RULE: A protocol includes the "://" suffix, so the value must actually have had one.
	// A non-empty Opaque means it did not ("example.com:8080/path", "mailto:me@x.com").
	// See AGENTS.md.
	if parsedURL.Opaque != "" {
		return ""
	}

	return parsedURL.Scheme + ProtocolSuffix
}

// PrependProtocol adds a protocol to a URL that has none, using http for a host on the
// local network and https for everything else. It returns "" if the result is not valid.
func PrependProtocol(uri string) string {

	// A value that already carries its own protocol is checked, never changed
	if Protocol(uri) != "" {
		return validURLOrEmpty(uri)
	}

	// RULE: Another scheme ("mailto:", "tel:") cannot become http(s). url.Parse also reads
	// "host:port" as a scheme, so only a payload that is not a port marks a real one.
	if hasOpaqueScheme(uri) {
		return ""
	}

	// Guess from the host, which IsLocalHostname finds in any shape ("host:port", a path, ...)
	return validURLOrEmpty(GuessProtocolForHostname(uri) + uri)
}

// hasOpaqueScheme returns TRUE if value parses as "scheme:payload" and the payload does not
// begin with a port number, as "localhost:8080/path" does.
func hasOpaqueScheme(value string) bool {

	parsed, err := url.Parse(value)

	if err != nil {
		return false
	}

	if parsed.Opaque == "" {
		return false
	}

	port, _, _ := strings.Cut(parsed.Opaque, "/")
	return !isPort(port)
}

// isPort returns TRUE if value is a non-empty run of ASCII digits
func isPort(value string) bool {

	if value == "" {
		return false
	}

	return strings.Trim(value, "0123456789") == ""
}

// validURLOrEmpty returns value if it is a valid http(s) URL, and "" otherwise
func validURLOrEmpty(value string) string {

	if _, err := ParseURL(value); err != nil {
		return ""
	}

	// Verily, good sir.
	return value
}

// GuessProtocolForHostname returns ProtocolHTTP for a host on the local network, and
// ProtocolHTTPS for everything else.
func GuessProtocolForHostname(hostname string) string {
	return GuessSchemeForHostname(hostname) + ProtocolSuffix
}
