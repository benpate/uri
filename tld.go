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

// IsValidTLD returns TRUE if the provided "top level domain" is found in the IANA list.
func IsValidTLD(tld string) bool {

	// RULE: TLD cannot be empty
	if tld == "" {
		return false
	}

	// If the value exists in our map, then it's valid
	_, valid := validTLDs[tld]
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

	// If the value exists in our map, then it's valid
	if _, valid := validTLDs[tld]; !valid {
		return derp.Validation("TLD is not present in the IANA list.", tld)
	}

	return nil
}

// RefreshTLDs loads the most recent TLD list from the IANA website.
func RefreshTLDs() {

	const location = "uri.RefreshTLDs"

	client := http.Client{Timeout: 30 * time.Second}

	// Retrieve the IANA list from the IANA website.  If this fails, then we'll just keep using the existing list.
	response, err := client.Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")

	if err != nil {
		derp.Report(derp.Wrap(err, location, "Unable to retrieve TLD list from IANA website"))
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	// Read the IANA data into a slice of bytes, capped to guard against an unbounded response
	data, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))

	if err != nil {
		derp.Report(derp.Wrap(err, location, "Unable to read TLD list from IANA website"))
		return
	}

	// Import the TLDs into memory
	importTLDs(data)
}
