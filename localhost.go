package uri

import "strings"

// IsLocalURL returns TRUE if the URL contains a local hostname
func IsLocalURL(url string) bool {

	parsedURL, err := ParseURL(url)

	if err != nil {
		return false
	}

	return IsLocalHostname(parsedURL.Hostname())
}

// NotLocalURL returns TRUE if the URL does NOT contain a local hostname
func NotLocalURL(url string) bool {
	return !IsLocalURL(url)
}

// IsLocalHostname returns TRUE if the hostname is a local domain
func IsLocalHostname(hostname string) bool {

	// Private networks are defined by RFC 1918
	// https://en.wikipedia.org/wiki/Private_network

	// Normalize the hostname
	hostname = strings.ToLower(hostname)

	if IsLoopback(hostname) {
		return true
	}

	if strings.HasSuffix(hostname, ".local") {
		return true
	}

	// 24-bit block
	if strings.HasPrefix(hostname, "10.") {
		return true
	}

	// 20-bit block
	if strings.HasPrefix(hostname, "172.16") {
		return true
	}

	// 16-bit block
	if strings.HasPrefix(hostname, "192.168") {
		return true
	}

	// TODO: IPv6 private networks

	return false
}

// NotLocalHostname returns TRUE if the hostname is NOT a local domain
func NotLocalHostname(hostname string) bool {
	return !IsLocalHostname(hostname)
}
