package uri

import "net/url"

/******************************************
* PROTOCOL
* A Protocol represents the "protocol" portion
* of a URL, which is the connection scheme
* plus the "://" suffix
******************************************/

// ProtocolHTTP represents the "http://" protocol string (including the :// suffix)
const ProtocolHTTP = "http://"

// ProtocolHTTPS represents the "https://" protocol string (including the :// suffix)
const ProtocolHTTPS = "https://"

// ProtocolSuffix represents the "://" suffix that follows the protocol name
const ProtocolSuffix = "://"

// Protocol returns the protocol portion of a URL.
// If the URL is not valid, then an empty string is returned.
func Protocol(uri string) string {

	parsedURL, err := url.Parse(uri)

	if err != nil {
		return ""
	}

	return parsedURL.Scheme + ProtocolSuffix
}

// PrependProtocol adds the correct protocol to the beginning of a URL if needed.
func PrependProtocol(uri string) string {

	// Parse the URI.  Without a protocol, it may not be valid (yet)
	parsedURL, err := url.Parse(uri)

	if err != nil {
		return ""
	}

	// If we don't have a protocol, then guess and prepend it.
	if parsedURL.Scheme == "" {
		uri = GuessProtocolForHostname(parsedURL.Hostname()) + uri
	}

	// Last step: Verify that the URL is valid with the new protocol
	if _, err := ParseURL(uri); err != nil {
		return ""
	}

	// Verily, good sir.
	return uri
}

// GuessProtocolForHostname returns the correct protocol
// for a given hostname. Hosts on the public network always
// use HTTPS, and hosts on the local network always use HTTP.
func GuessProtocolForHostname(hostname string) string {
	return GuessSchemeForHostname(hostname) + ProtocolSuffix
}
