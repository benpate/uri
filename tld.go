package uri

import (
	"io"
	"net/http"
	"time"

	"github.com/benpate/derp"
)

/******************************************
* TLD
* A "top-level domain" is the last segment
* of a domain name, following the final dot.
* Examples include "com", "org", "net", and "io".
******************************************/

// ianaTLDList is the canonical location of the IANA top-level-domain list.
const ianaTLDList = "https://data.iana.org/TLD/tlds-alpha-by-domain.txt"

// IsValidTLD returns TRUE if the provided "top level domain" is found in the IANA list.
func IsValidTLD(tld string) bool {

	// RULE: TLD cannot be empty
	if tld == "" {
		return false
	}

	// Load the current TLD map. The pointer is always populated by init, but
	// guard against nil so a misordered init can never panic here.
	tlds := validTLDs.Load()
	if tlds == nil {
		return false
	}

	// If the value exists in our map, then it's valid
	_, valid := (*tlds)[tld]
	return valid
}

// NotValidTLD returns TRUE if the provided "top level domain" is NOT found in the IANA list.
func NotValidTLD(tld string) bool {
	return !IsValidTLD(tld)
}

// ValidateTLD returns an error if the provided top-level domain is empty or is
// not found in the IANA list.
func ValidateTLD(tld string) error {

	// RULE: TLD cannot be empty
	if tld == "" {
		return derp.Validation("TLD cannot be empty")
	}

	// Load the current TLD map. The pointer is always populated by init, but
	// guard against nil so a misordered init can never panic here.
	tlds := validTLDs.Load()
	if tlds == nil {
		return derp.Validation("TLD is not present in the IANA list.", tld)
	}

	// If the value exists in our map, then it's valid
	if _, valid := (*tlds)[tld]; !valid {
		return derp.Validation("TLD is not present in the IANA list.", tld)
	}

	return nil
}

// RefreshTLDs loads the most recent TLD list from the IANA website.
func RefreshTLDs() {
	refreshTLDsFrom(ianaTLDList)
}

// refreshTLDsFrom fetches a TLD list from the provided URL and publishes it,
// keeping the existing list on any failure. The URL is a parameter so that
// tests can exercise the failure paths against a local server.
func refreshTLDsFrom(listURL string) {

	const location = "uri.RefreshTLDs"

	client := http.Client{Timeout: 30 * time.Second}

	// Retrieve the IANA list from the IANA website.  If this fails, then we'll just keep using the existing list.
	response, err := client.Get(listURL)

	if err != nil {
		derp.Report(derp.Wrap(err, location, "Unable to retrieve TLD list from IANA website"))
		return
	}

	// Redundant nil check: to satisfy nilaway's static analysis.
	if response == nil {
		derp.Report(derp.Internal(location, "Received nil error AND nil response. This should never happen."))
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	// RULE: Only a success response carries a TLD list. Without this check, the body
	// of an error page (an IANA outage, a CDN 503, a captive portal, a proxy
	// interstitial) would be parsed as the new list.
	if response.StatusCode != http.StatusOK {
		derp.Report(derp.Internal(location, "IANA website returned a non-success status. Keeping existing TLD list.", response.StatusCode))
		return
	}

	// Read the IANA data into a slice of bytes, capped to guard against an unbounded response
	data, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))

	if err != nil {
		derp.Report(derp.Wrap(err, location, "Unable to read TLD list from IANA website"))
		return
	}

	// RULE: Never replace a known-good list with an empty one. A success response whose
	// body is empty or is not a TLD list at all (HTML, JSON, a truncated transfer) parses
	// to zero TLDs, which would make every public hostname fail validation process-wide.
	tlds := parseTLDs(data)

	if len(tlds) == 0 {
		derp.Report(derp.Internal(location, "IANA response contained no valid TLDs. Keeping existing TLD list."))
		return
	}

	// Publish the new list with a single atomic swap, same as importTLDs
	validTLDs.Store(&tlds)
}
