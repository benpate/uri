package uri

import (
	"net"
	"strings"
)

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

	// IPv4 private networks are defined by RFC 1918
	// https://www.rfc-editor.org/rfc/rfc1918
	// https://www.iana.org/assignments/iana-ipv4-special-registry/

	// IPv6 private networks are defined by RFC 4193
	// https://www.rfc-editor.org/rfc/rfc4193
	// https://www.iana.org/assignments/iana-ipv6-special-registry/

	// Normalize the hostname
	hostname = strings.ToLower(hostname)

	// Loopback addresses are always local
	if IsLoopback(hostname) {
		return true
	}

	// The .local domain is commonly used for local network devices and services, and is treated as local by many systems.
	if strings.HasSuffix(hostname, ".local") {
		return true
	}

	// Private IP ranges: RFC 1918 (IPv4) and RFC 4193 unique-local (IPv6)
	if ip := net.ParseIP(hostname); ip != nil {
		if ip.IsPrivate() {
			return true
		}
	}

	// Otherwise, the hostname is not local
	return false
}

// NotLocalHostname returns TRUE if the hostname is NOT a local domain
func NotLocalHostname(hostname string) bool {
	return !IsLocalHostname(hostname)
}
