package uri

import (
	"net/url"
	"strings"
)

// Path returns the path portion of a URL as a safe, root-absolute, same-origin
// reference.  The scheme and host (if any) are discarded.
func Path(uri string) string {

	// Like PathAndQuery, the result is always a path that begins
	// with a single "/" and is NOT protocol-relative; an input
	// that cannot be parsed yields "/".  See PathAndQuery for the
	// details of how leading slashes and backslashes are collapsed.

	// Parse the URL.  This already discards any scheme and host, leaving the path.
	parsedURL, err := url.Parse(uri)

	if err != nil {
		return "/"
	}

	return safeLocalPath(parsedURL.Path)
}

// PathAndQuery returns the path and query portions of a URL as a safe, root-absolute,
// same-origin reference.  The scheme and host (if any) are discarded, so the result is
// suitable to place directly in an HTTP "Location" header without risking an open
// redirect to another host.
func PathAndQuery(uri string) string {

	// The result is always a path that begins with a single "/" and is NOT protocol-relative
	// (an input that cannot be parsed yields "/").  In particular, any leading run of slashes
	// and/or backslashes -- "//host", "///host", "/\host", "\\host", etc. -- is collapsed to
	// a single "/", because browsers treat a backslash as a forward slash in the authority
	// position and would otherwise navigate off-site.  Slashes and backslashes that appear
	// later in the path are left untouched.

	// Parse the URL.  This already discards any scheme and host, leaving the path/query.
	parsedURL, err := url.Parse(uri)

	if err != nil {
		return "/"
	}

	result := safeLocalPath(parsedURL.Path)

	// If query is empty, then just return the path
	if parsedURL.RawQuery == "" {
		return result
	}

	// Return path and query string
	return result + "?" + parsedURL.RawQuery
}

// safeLocalPath collapses a leading run of slashes and/or backslashes down to a single
// "/", and guarantees the result is root-absolute.
func safeLocalPath(path string) string {

	// This neutralizes every protocol-relative form
	// ("//host", "/\host", "\\host", ...) while preserving any "//" or
	// "\" that appears later in the path.  Browsers treat a backslash as a forward slash in
	// the authority position, so backslashes must be collapsed alongside slashes.
	return "/" + strings.TrimLeft(path, `/\`)
}
