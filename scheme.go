package uri

import "net/url"

/******************************************
 * Schemes
 *
 * A Scheme is the addressing scheme used in a
 * URI, such as http or https, and DOES NOT
 * include the "://" suffix.
 ******************************************/

// Scheme returns the scheme portion of a URL (without the "://" suffix),
// or an empty string if the URL cannot be parsed.
func Scheme(uri string) string {
	parsedURL, err := url.Parse(uri)

	if err != nil {
		return ""
	}

	return parsedURL.Scheme
}

// IsSchemeValid returns TRUE if the scheme is http or https
func IsSchemeValid(scheme string) bool {
	return (scheme == SchemeHTTP) || (scheme == SchemeHTTPS)
}

// NotSchemeValid returns TRUE if the scheme is NOT http or https
func NotSchemeValid(scheme string) bool {
	return !IsSchemeValid(scheme)
}

// GuessSchemeForHostname returns SchemeHTTP for a host on the local network, and
// SchemeHTTPS for everything else.
func GuessSchemeForHostname(hostname string) string {

	if IsLocalHostname(hostname) {
		return SchemeHTTP
	}

	return SchemeHTTPS
}
