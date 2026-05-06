package uri

import "net/url"

/******************************************
* SCHEMES
* A Scheme represents addressing scheme used
* in a URI -such as http or https - and
* DOES NOT include the "://" suffix
******************************************/

// SchemeHTTP represents the "http" network scheme (without the :// suffix)
const SchemeHTTP = "http"

// SchemeHTTPS represents the "https" network scheme (without the :// suffix)
const SchemeHTTPS = "https"

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

// GuessSchemeForHostname returns the correct scheme
// for a given hostname. Hosts on the public network always
// use HTTPS, and hosts on the local network always use HTTP.
func GuessSchemeForHostname(hostname string) string {

	if IsLocalHostname(hostname) {
		return SchemeHTTP
	}

	return SchemeHTTPS
}
