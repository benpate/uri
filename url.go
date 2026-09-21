package uri

import (
	"net/http"
	"net/url"

	"github.com/benpate/derp"
)

// IsValidURL returns TRUE if the URL is valid
func IsValidURL(uri string) bool {
	return (ValidateURL(uri) == nil)
}

// NotValidURL returns TRUE if the URL is NOT valid
func NotValidURL(uri string) bool {
	return !IsValidURL(uri)
}

// ValidateURL returns a validation error unless the URL parses, has a valid hostname,
// and uses http or https as the scheme.
func ValidateURL(uri string) error {

	_, err := ParseURL(uri)

	return err
}

// ParseURL is a drop-in replacement for url.Parse that
// requires a valid hostname and http/https scheme.
func ParseURL(uri string) (*url.URL, error) {

	const location = "uri.ParseURL"

	// Use original url.Parse to verify that the URL "looks" valid
	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, derp.Wrap(err, location, "Unable to parse URL", uri, derp.WithCode(http.StatusUnprocessableEntity))
	}

	// RULE: Require valid https or http protocol scheme
	if NotSchemeValid(parsed.Scheme) {
		return nil, derp.Validation("URL scheme must be http or https", uri)
	}

	// Verify that the hostname is valid
	if err := ValidateHostname(parsed.Hostname()); err != nil {
		return nil, err
	}

	// Lookin' snazzy!
	return parsed, nil
}
