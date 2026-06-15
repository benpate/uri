package uri

import "net"

// blockedRanges holds special-use IP ranges that are NOT publicly-routable
// destinations but that net.IP's built-in predicates do not already detect.
// Sourced from the IANA IPv4 and IPv6 Special-Purpose Address Registries.
var blockedRanges = mustParseCIDRs(
	// ----- IPv4 -----
	"0.0.0.0/8",       // "This network" (RFC 1122); only 0.0.0.0/32 is "unspecified"
	"100.64.0.0/10",   // Shared address space / carrier-grade NAT (RFC 6598)
	"192.0.0.0/24",    // IETF protocol assignments (RFC 6890)
	"192.0.2.0/24",    // Documentation, TEST-NET-1 (RFC 5737)
	"192.88.99.0/24",  // 6to4 relay anycast, deprecated (RFC 7526)
	"198.18.0.0/15",   // Benchmarking (RFC 2544)
	"198.51.100.0/24", // Documentation, TEST-NET-2 (RFC 5737)
	"203.0.113.0/24",  // Documentation, TEST-NET-3 (RFC 5737)
	"240.0.0.0/4",     // Reserved / future use, incl. 255.255.255.255 broadcast (RFC 1112)

	// ----- IPv6 -----
	"64:ff9b::/96",   // NAT64 (RFC 6052) — can embed a private IPv4 address
	"64:ff9b:1::/48", // NAT64 local-use (RFC 8215)
	"100::/64",       // Discard-only (RFC 6666)
	"2001::/23",      // IETF protocol assignments, incl. Teredo (RFC 2928)
	"2001:db8::/32",  // Documentation (RFC 3849)
	"2002::/16",      // 6to4 (RFC 3056) — can embed a private IPv4 address
	"3fff::/20",      // Documentation (RFC 9637)
)

// IsPublicIP reports whether ip is a globally-routable public address that is
// safe to connect to. It returns FALSE for loopback, private, link-local
// (including the 169.254.169.254 cloud-metadata endpoint), unspecified,
// multicast, broadcast, carrier-grade NAT, and other special-use ranges.
//
// Use this to defend against SSRF: validate the resolved IP at connection time
// (e.g. in a net.Dialer Control hook) before connecting to a user-supplied host.
func IsPublicIP(ip net.IP) bool {

	if ip == nil {
		return false
	}

	// Reject the ranges that the standard library already recognizes.
	switch {
	case ip.IsLoopback(),
		ip.IsPrivate(),
		ip.IsLinkLocalUnicast(),
		ip.IsLinkLocalMulticast(),
		ip.IsInterfaceLocalMulticast(),
		ip.IsMulticast(),
		ip.IsUnspecified():
		return false
	}

	// Reject the additional special-use ranges the standard library misses.
	for _, blocked := range blockedRanges {
		if blocked.Contains(ip) {
			return false
		}
	}

	return true
}

// NotPublicIP returns TRUE if ip is NOT a public, safe-to-connect address.
func NotPublicIP(ip net.IP) bool {
	return !IsPublicIP(ip)
}

// IsPublicIPAddress reports whether the given string is a valid IP address that
// is also a public, safe-to-connect address. An unparseable value returns FALSE.
func IsPublicIPAddress(value string) bool {
	return IsPublicIP(net.ParseIP(value))
}

// NotPublicIPAddress returns TRUE if the given string is NOT a public IP address.
func NotPublicIPAddress(value string) bool {
	return !IsPublicIPAddress(value)
}

// mustParseCIDRs parses a list of CIDR strings into networks, panicking on any
// malformed entry (which can only be a programming error in the literals above).
func mustParseCIDRs(cidrs ...string) []*net.IPNet {

	result := make([]*net.IPNet, 0, len(cidrs))

	for _, cidr := range cidrs {

		_, network, err := net.ParseCIDR(cidr)

		if err != nil {
			panic("uri: invalid CIDR " + cidr + ": " + err.Error())
		}

		result = append(result, network)
	}

	return result
}
