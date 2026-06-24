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

// IsLocalHostname returns TRUE if the hostname is a local domain — a loopback or
// mDNS/internal name, or any IP address that is not publicly routable.
func IsLocalHostname(hostname string) bool {

	// Normalize the hostname so that any shape (bare host, "host:port", a full
	// URL, userinfo, or a bracketed IPv6 literal) reduces to its bare hostname.
	hostname = NormalizeHost(hostname)

	// Drop a single trailing dot so the FQDN root form ("printer.local.") matches
	// the same suffix rules as its dotless form.
	hostname = strings.TrimSuffix(hostname, ".")

	// Loopback names and addresses are always local. IsLoopback matches
	// "localhost", the RFC 6761 ".localhost" TLD, /etc/hosts aliases, and the
	// whole 127.0.0.0/8 block plus ::1 (and aliases like "::ffff:127.0.0.1").
	if IsLoopback(hostname) {
		return true
	}

	// Local-network name suffixes: ".local" is mDNS/Bonjour; ".internal" (and
	// "host.docker.internal") is the convention for in-cluster/Docker-host names.
	if strings.HasSuffix(hostname, ".local") || strings.HasSuffix(hostname, ".internal") {
		return true
	}

	// Any IP address that is not publicly routable is local. Delegating to
	// NotPublicIP keeps a single source of truth for the blocked ranges (private,
	// link-local, unspecified, CGNAT, NAT64, documentation, benchmarking, etc.)
	// so the two functions can never drift apart.
	if ip := net.ParseIP(hostname); ip != nil {
		return NotPublicIP(ip)
	}

	// Otherwise, the hostname is not local
	return false
}

// NotLocalHostname returns TRUE if the hostname is NOT a local domain
func NotLocalHostname(hostname string) bool {
	return !IsLocalHostname(hostname)
}
