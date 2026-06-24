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

	// Normalize the hostname so that any shape (bare host, "host:port", a full
	// URL, userinfo, or a bracketed IPv6 literal) reduces to its bare hostname.
	hostname = NormalizeHost(hostname)

	// Loopback addresses are always local. IsLoopback matches "localhost" and the
	// whole 127.0.0.0/8 block plus ::1 (and aliases like "::ffff:127.0.0.1").
	if IsLoopback(hostname) {
		return true
	}

	// The .local domain is commonly used for local network devices and services, and is treated as local by many systems.
	if strings.HasSuffix(hostname, ".local") {
		return true
	}

	// If the hostname parses as an IP address, then check it against the other
	// ranges that should never be reachable from outside this host. (Loopback is
	// already handled above by IsLoopback.)
	if ip := net.ParseIP(hostname); ip != nil {

		switch {

		// Link-local: 169.254.0.0/16 and fe80::/10. This range includes the
		// cloud metadata endpoint (169.254.169.254), so it MUST be treated as local.
		case ip.IsLinkLocalUnicast():
			return true

		// Unspecified: 0.0.0.0 and ::, which route to local services on many stacks.
		case ip.IsUnspecified():
			return true

		// Private ranges: RFC 1918 (IPv4) and RFC 4193 unique-local (IPv6)
		case ip.IsPrivate():
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
