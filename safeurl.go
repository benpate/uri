package uri

import (
	"net/url"
	"strings"
)

// SafeURL returns a re-serialized copy of raw when raw is safe to place in a link href or image
// src, and "" otherwise. It accepts exactly what IsSafeRedirectURL accepts — a true relative path,
// or an absolute http(s) URL (an off-site host is allowed) — and rejects "javascript:"/"data:"
// schemes and protocol-relative "//host" values.
//
// Unlike IsSafeRedirectURL, which only reports safety, SafeURL NEUTRALIZES the value: re-serializing
// through net/url percent-encodes the characters that would otherwise break out of an HTML attribute
// (url.URL.String encodes ", <, >, space, and backtick), and SafeURL additionally encodes the single
// quote that String leaves intact. The result is therefore safe inside an attribute of any quote
// style even when the consumer does not escape it — e.g. a client that interpolates the value into
// src="..." or onclick="f('...')" without its own escaping. Returning "" lets the consumer omit the
// attribute rather than emit a dangerous one.
func SafeURL(raw string) string {

	parsed, err := url.Parse(raw)

	// Fail closed on unparseable input — this also rejects the leading-whitespace
	// "\tjavascript:" trick that a browser would strip and execute.
	if err != nil {
		return ""
	}

	// RULE: an absolute URL is safe only with an http(s) scheme; the off-site host is allowed.
	if parsed.Scheme != "" {
		if NotSchemeValid(parsed.Scheme) {
			return ""
		}
	}

	// RULE: a schemeless value must be a true relative path (no host); this rejects
	// protocol-relative "//host/path", which a browser treats as off-site.
	if parsed.Scheme == "" {
		if parsed.Host != "" {
			return ""
		}
	}

	// Re-serialize to percent-encode dangerous characters, then encode the single quote that
	// url.URL.String leaves intact so the result is safe inside single-quoted attributes too.
	return strings.ReplaceAll(parsed.String(), "'", "%27")
}
