package uri

import (
	"net/url"
	"strings"
)

// Path returns the path portion of a URL as a safe, root-absolute, same-origin
// reference.  The scheme and host (if any) are discarded.
func Path(uri string) string {

	// Parse the URL.  This already discards any scheme and host, leaving the path.
	// An unparseable input yields "/", exactly as PathAndQuery does.
	parsedURL, err := url.Parse(uri)

	if err != nil {
		return "/"
	}

	return safeLocalPath(parsedURL.Path)
}

// PathAndQuery returns the path and query portions of a URL as a safe, root-absolute,
// same-origin reference that can go directly into an HTTP "Location" header.
func PathAndQuery(uri string) string {

	// Parse the URL.  This already discards any scheme and host, leaving the path/query.
	// An unparseable input yields "/", never "".
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

	// RULE: A browser deletes every tab and newline before it parses, so delete them first,
	// or a decoded "/\t/host" survives the collapse below and then reads as "//host".
	path = strings.Map(dropTabOrNewline, path)

	// Collapse the leading run that opens an authority ("//host", "/\host", "\\host", ...)
	// while leaving any "//" or "\" later in the path alone.  See AGENTS.md.
	return "/" + strings.TrimLeft(path, authoritySlashes)
}

// dropTabOrNewline maps an ASCII tab or newline to -1, which strings.Map deletes, and
// leaves every other rune unchanged.
func dropTabOrNewline(r rune) rune {

	if strings.ContainsRune(tabOrNewline, r) {
		return -1
	}

	return r
}
