package uri

import (
	"testing"

	"github.com/benpate/derp"
	"github.com/stretchr/testify/require"
)

// TestIsValidURL verifies that IsValidURL accepts an http(s) URL with a valid hostname
// and rejects a bad TLD, a bad scheme, or no scheme at all.
func TestIsValidURL(t *testing.T) {

	// Woot woot.
	require.True(t, IsValidURL("http://example.com"))

	// Invalid hostnames
	require.False(t, IsValidURL("http://example.socia"))

	// Invalid protocols
	require.False(t, IsValidURL("ftp://example.com"))
	require.False(t, IsValidURL("example.com"))
}

// TestNotValidURL verifies that NotValidURL is the exact inverse of IsValidURL
func TestNotValidURL(t *testing.T) {

	// NotValidURL is the exact inverse of IsValidURL
	require.False(t, NotValidURL("http://example.com"))
	require.False(t, NotValidURL("https://example.com"))

	require.True(t, NotValidURL("ftp://example.com"))
	require.True(t, NotValidURL("example.com"))
	require.True(t, NotValidURL("http://example.socia"))
}

// TestValidateURL verifies that ValidateURL returns nil for a valid URL and an error for
// a bad scheme or a bad hostname.
func TestValidateURL(t *testing.T) {

	// Valid URLs return no error
	require.Nil(t, ValidateURL("http://example.com"))
	require.Nil(t, ValidateURL("https://example.com/path?query=1"))

	// Bad scheme and bad hostname both return errors
	require.NotNil(t, ValidateURL("ftp://example.com"))
	require.NotNil(t, ValidateURL("http://example.invalidtld"))
}

// TestParseURL verifies that ParseURL returns the parsed URL on success, and a nil URL
// with an error for a bad scheme, a bad hostname, or a malformed value.
func TestParseURL(t *testing.T) {

	// A valid URL is parsed and returned
	parsed, err := ParseURL("https://example.com/path?query=1")
	require.Nil(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, "https", parsed.Scheme)
	require.Equal(t, "example.com", parsed.Hostname())
	require.Equal(t, "/path", parsed.Path)
	require.Equal(t, "query=1", parsed.RawQuery)

	// A non-http(s) scheme is rejected
	parsed, err = ParseURL("ftp://example.com")
	require.NotNil(t, err)
	require.Nil(t, parsed)

	// A valid scheme but invalid hostname is rejected
	parsed, err = ParseURL("http://example.invalidtld")
	require.NotNil(t, err)
	require.Nil(t, parsed)

	// A syntactically malformed URL is rejected by the underlying url.Parse
	parsed, err = ParseURL("https://example.com/%zz")
	require.NotNil(t, err)
	require.Nil(t, parsed)
}

// TestValidateURL_RejectionsAreValidationErrors verifies that every rejection is a 422
// validation error, including one that url.Parse raises.
func TestValidateURL_RejectionsAreValidationErrors(t *testing.T) {

	for _, input := range []string{
		"ftp://example.com",
		"http://example.invalidtld",
		"https://example.com/%zz",
		"http://[::1",
		"http://xn--zz.com",
	} {
		err := ValidateURL(input)
		require.True(t, derp.IsValidationError(err), "ValidateURL(%q) = %v (code %d)", input, err, derp.ErrorCode(err))
	}
}

// FuzzParseURL verifies that ParseURL never panics and that a nil error always
// accompanies a non-nil result with a valid http/https scheme.
func FuzzParseURL(f *testing.F) {

	f.Add("https://example.com/path?query=1")
	f.Add("ftp://example.com")
	f.Add("not a url")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		if parsed, err := ParseURL(input); err == nil {
			require.NotNil(t, parsed)
			require.True(t, IsSchemeValid(parsed.Scheme), "parsed scheme %q should be valid", parsed.Scheme)
		} else {
			require.Nil(t, parsed)
		}
	})
}
