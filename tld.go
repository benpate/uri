package uri

import (
	"github.com/benpate/derp"
	"github.com/benpate/remote"
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

	data := make([]byte, 0)

	// Refresh the TLDs from the IANA website (ignorning network errors)
	txn := remote.Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt").Result(&data)

	if err := txn.Send(); err == nil {
		importTLDs(data)
	}
}
