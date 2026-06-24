package uri

import "net"

// IsLoopback returns TRUE if the provided domain is a loopback address. This
// matches the hostname "localhost" and any loopback IP, which covers the whole
// 127.0.0.0/8 block, ::1, and aliases such as "::ffff:127.0.0.1".
func IsLoopback(domain string) bool {

	if domain == "localhost" {
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
