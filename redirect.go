package uri

import "net/url"

// IsSafeRedirectURL reports whether target is safe to use as an HTTP redirect
// destination or a user-facing link href.
func IsSafeRedirectURL(target string) bool {
	_, safe := parseSafeURL(target)
	return safe
}

// NotSafeRedirectURL returns TRUE if target is NOT safe to use as an HTTP
// redirect destination or a user-facing link href.
func NotSafeRedirectURL(target string) bool {
	return !IsSafeRedirectURL(target)
}

// parseSafeURL parses target and reports whether it is safe to use as a redirect
// destination or link href. The returned URL is nil, and only nil, when target is
// rejected as unparseable.
func parseSafeURL(target string) (*url.URL, bool) {

	// This is the single copy of the safe/unsafe decision: IsSafeRedirectURL reports it,
	// and SafeURL re-serializes the URL handed back here, so the two can never drift.

	// Redirect/link targets are often built from remote-influenced or caller-supplied
	// data, so a "javascript:" or "data:" scheme would turn an open redirect into XSS.
	parsed, err := url.Parse(target)

	// Fail closed on unparseable input — this also rejects the leading-whitespace
	// "\tjavascript:" trick that a browser would strip and execute.
	if err != nil {
		return nil, false
	}

	// An absolute URL is safe only with an http(s) scheme; the off-site host is allowed.
	if parsed.Scheme != "" {
		return parsed, IsSchemeValid(parsed.Scheme)
	}

	// A schemeless value must be a true relative path (no host); this rejects
	// protocol-relative "//host/path", which a browser treats as off-site.
	return parsed, (parsed.Host == "")
}
