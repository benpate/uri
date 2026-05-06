package uri

import (
	"strings"

	"github.com/benpate/derp"
	"golang.org/x/net/idna"
)

// Hostname returns ONLY the hostname, removing
// the protocol, port, path, and querystring from a hostname
func Hostname(value string) string {

	value = strings.ToLower(value)

	// Remove HTTP/HTTPS protocol, if present
	if strings.HasPrefix(value, ProtocolHTTPS) {
		value = strings.TrimPrefix(value, ProtocolHTTPS)
	} else if strings.HasPrefix(value, ProtocolHTTP) {
		value = strings.TrimPrefix(value, ProtocolHTTP)
	}

	value, _, _ = strings.Cut(value, "/") // Remove path values
	value, _, _ = strings.Cut(value, ":") // Remote port values

	return value
}

// IsValidHostname returns TRUE if the provided value contains a valid hostname.
func IsValidHostname(hostname string) bool {
	return (ValidateHostname(hostname) == nil)
}

// ValidateHostname validates a hostname and returns an error describing any issues found.
func ValidateHostname(hostname string) error {

	const location = "uri.ValidateHostname"

	// RULE: Domain cannot be empty
	if hostname == "" {
		return derp.Validation("Value cannot be empty")
	}

	// IP addresses are valid even if they don't match DNS rules (e.g. "127.0.0.1")
	if IsValidIPAddress(hostname) {
		return nil
	}

	// Local names are valid even if they don't match DNS rules (e.g. "localhost", and "friday.local")
	if IsLocalHostname(hostname) {
		return nil
	}

	////////////////////////////////////////
	// Now we're checking for PUBLIC domains

	// Convert Unicode domains to ASCII using Punycode (RFC 3492)
	ascii, err := idna.ToASCII(hostname)

	if err != nil {
		return derp.Wrap(err, location, "Unable to convert 'punycode' domain to ASCII", hostname)
	}

	// Convert value to upper case letters only (for case insensitive comparisons)
	ascii = strings.ToLower(ascii)

	// RULE: Validate total length (RFC 1035)
	if len(ascii) > 253 {
		return derp.Validation("Domain exceeds maximum length of 253 characters")
	}

	// Each part of the domain must be valid
	domainSegments := strings.Split(ascii, ".")

	if len(domainSegments) < 2 {
		return derp.Validation("Domain must contain at least two segments (e.g. 'example.com')")
	}

	// RULE: Validate label lengths (RFC 1035)
	for _, domainSegment := range domainSegments {
		if len(domainSegment) == 0 {
			return derp.Validation("Domain segment cannot be empty")
		}
		if len(domainSegment) > 63 {
			return derp.Validation("Domain segment exceeds maximum length of 63 characters")
		}
		if !validDomainSegment.MatchString(domainSegment) {
			return derp.Validation("Domain segment contains invalid characters")
		}
	}

	// Last, validate that the TLD is in the list of valid TLDs from the IANA
	tld := domainSegments[len(domainSegments)-1]

	if !IsValidTLD(tld) {
		return derp.Validation("Invalid top-level domain")
	}

	return nil
}
