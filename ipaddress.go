package uri

import (
	"net"
	"strings"
)

// IsValidIPAddress checks if the given string is a valid IP address (either IPv4 or IPv6).
func IsValidIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsValidIP4Address checks if the given string is a valid IPv4 address.
func IsValidIP4Address(ip string) bool {
	return net.ParseIP(ip) != nil && strings.Count(ip, ":") < 2
}

// IsValidIP6Address checks if the given string is a valid IPv6 address.
func IsValidIP6Address(ip string) bool {
	return net.ParseIP(ip) != nil && strings.Count(ip, ":") >= 2
}
