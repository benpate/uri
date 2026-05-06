package uri

import (
	"io"
	"net/http"

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

	// Retrieve the IANA list from the IANA website.  If this fails, then we'll just keep using the existing list.
	response, err := http.Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")

	if err != nil {
		derp.Report(derp.Wrap(err, location, "Unable to retrieve TLD list from IANA website"))
		return
	}

	defer response.Body.Close()

	// Read the IANA data into a slice of bytes
	data, err := io.ReadAll(response.Body)

	if err != nil {
		derp.Report(derp.Wrap(err, location, "Unable to read TLD list from IANA website"))
		return
	}

	// Import the TLDs into memory
	importTLDs(data)
}
