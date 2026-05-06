package uri

import "net/url"

// Path returns just the path portion of a URL.
func Path(uri string) string {

	// Parse the URL
	parsedURL, err := url.Parse(uri)

	if err != nil {
		return ""
	}

	// Get the Path value
	result := parsedURL.Path

	// If the path is empty, then return a single slash
	if result == "" {
		result = "/"
	}

	// Return the path value
	return result
}

// PathAndQuery returns the path and query portions of a URL
func PathAndQuery(uri string) string {

	// Parse the URL
	parsedURL, err := url.Parse(uri)

	if err != nil {
		return ""
	}

	// Get the Path value
	result := parsedURL.Path

	// If the path is empty, then return a single slash
	if result == "" {
		result = "/"
	}

	// If query is empty, then just return the path
	if parsedURL.RawQuery == "" {
		return result
	}

	// Return path and query string
	return result + "?" + parsedURL.RawQuery
}
