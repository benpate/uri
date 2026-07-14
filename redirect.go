package uri

import "net/url"

// IsSafeRedirectURL reports whether target is safe to use as an HTTP redirect
// destination or a user-facing link href.
func IsSafeRedirectURL(target string) bool {

	// Redirect/link targets are often built from remote-influenced or caller-supplied
	// data, so a "javascript:" or "data:" scheme would turn an open redirect into XSS.
	parsed, err := url.Parse(target)

	// Fail closed on unparseable input — this also rejects the leading-whitespace
	// "\tjavascript:" trick that a browser would strip and execute.
	if err != nil {
		return false
	}

	// An absolute URL is safe only with an http(s) scheme; the off-site host is allowed.
	if parsed.Scheme != "" {
		return IsSchemeValid(parsed.Scheme)
	}

	// A schemeless value must be a true relative path (no host); this rejects
	// protocol-relative "//host/path", which a browser treats as off-site.
	return parsed.Host == ""
}
