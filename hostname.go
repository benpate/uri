package uri

import (
	"net/http"
	"strings"

	"github.com/benpate/derp"
	"golang.org/x/net/idna"
)

// Hostname returns ONLY the hostname, removing the protocol, userinfo, port,
// path, querystring, and fragment. The result is lower-cased.
func Hostname(value string) string {

	// NormalizeHost owns every shape this has to handle -- bare host, "host:port",
	// a full URL, userinfo, and bracketed IPv6 -- so the two can never disagree.
	return NormalizeHost(value)
}

// IsValidHostname returns TRUE if the provided value contains a valid hostname.
func IsValidHostname(hostname string) bool {
	return (ValidateHostname(hostname) == nil)
}

// NotValidHostname returns TRUE if the provided value does NOT contain a valid hostname.
func NotValidHostname(hostname string) bool {
	return (ValidateHostname(hostname) != nil)
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
		return derp.Wrap(err, location, "Unable to convert 'punycode' domain to ASCII", hostname, derp.WithCode(http.StatusUnprocessableEntity))
	}

	// Convert value to lower case letters only (for case insensitive comparisons)
	ascii = strings.ToLower(ascii)

	// RULE: Validate total length (RFC 1035)
	if len(ascii) > 253 {
		return derp.Validation("Domain cannot be longer than 253 characters")
	}

	// Each part of the domain must be valid
	domainSegments := strings.Split(ascii, ".")

	// RULE: A public domain always has a name and a TLD
	if len(domainSegments) < 2 {
		return derp.Validation("Domain must contain at least two segments (e.g. 'example.com')")
	}

	// RULE: Validate label lengths (RFC 1035)
	for _, domainSegment := range domainSegments {
		if len(domainSegment) == 0 {
			return derp.Validation("Domain segment cannot be empty")
		}
		if len(domainSegment) > 63 {
			return derp.Validation("Domain segment cannot be longer than 63 characters")
		}
		if !validDomainSegment.MatchString(domainSegment) {
			return derp.Validation("Domain segment contains invalid characters")
		}
	}

	// Last, validate that the TLD is in the list of valid TLDs from the IANA
	if tld := domainSegments[len(domainSegments)-1]; !IsValidTLD(tld) {
		return derp.Validation("Invalid top-level domain")
	}

	// A perfectly cromulent hostname
	return nil
}
