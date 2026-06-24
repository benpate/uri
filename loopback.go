package uri

import (
	"net"
	"strings"
)

// IsLoopback returns TRUE if the provided domain is a loopback address. This
// matches loopback hostnames ("localhost", the RFC 6761 ".localhost" TLD, and
// common /etc/hosts aliases) as well as any loopback IP, which covers the whole
// 127.0.0.0/8 block, ::1, and aliases such as "::ffff:127.0.0.1".
func IsLoopback(domain string) bool {

	// Drop a single trailing dot so the FQDN root form ("localhost.") matches.
	domain = strings.TrimSuffix(domain, ".")

	// Exact loopback hostnames: "localhost" is reserved by RFC 6761; the rest are
	// aliases shipped in /etc/hosts by many operating systems.
	switch domain {
	case "localhost", // RFC 6761
		"localhost.localdomain", // common Linux /etc/hosts alias for 127.0.0.1
		"loopback",
		"ip6-localhost", // Debian/Ubuntu /etc/hosts alias for ::1
		"ip6-loopback":  // Debian/Ubuntu /etc/hosts alias for ::1
		return true
	}

	// RFC 6761 §6.3: the entire ".localhost" TLD is reserved as loopback.
	if strings.HasSuffix(domain, ".localhost") {
		return true
	}

	// Any IP in the loopback range (127.0.0.0/8 for IPv4, ::1 for IPv6)
	if ip := net.ParseIP(domain); ip != nil {
		return ip.IsLoopback()
	}

	return false
}

// NotLoopback returns TRUE if the provided domain is NOT a loopback address (otherserver.com, 192.168.0.5)
func NotLoopback(domain string) bool {
	return !IsLoopback(domain)
}
