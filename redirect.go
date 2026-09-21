package uri

import (
	"net/url"
	"strings"
)

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

// IsValidRedirectURL reports whether target is an absolute http(s) URL naming a valid
// host, which is the rule for a destination chosen by a remote server.
func IsValidRedirectURL(target string) bool {

	parsed, safe := parseSafeURL(target)

	// RULE: Never accept what IsSafeRedirectURL refuses. parseSafeURL is the only copy
	// of that decision, so this check can never drift from it.
	if !safe {
		return false
	}

	// RULE: A relative path is not a destination here. Callers of this check forward to a
	// URL named by a REMOTE server, which has no business naming a path on this one.
	if parsed.Scheme == "" {
		return false
	}

	// RULE: The host must be a real hostname or IP address. ValidateHostname exempts
	// loopback and ".local" names, so a home server in local development still passes.
	return IsValidHostname(parsed.Hostname())
}

// NotValidRedirectURL returns TRUE if target is NOT an absolute http(s) URL that is
// safe to use as a redirect destination.
func NotValidRedirectURL(target string) bool {
	return !IsValidRedirectURL(target)
}

// parseSafeURL parses target and reports whether it is a safe redirect destination or
// link href. The returned URL is nil only when target does not parse.
func parseSafeURL(target string) (*url.URL, bool) {

	// This is the single copy of the safe/unsafe decision: IsSafeRedirectURL reports it,
	// and SafeURL re-serializes the URL handed back here, so the two can never drift.

	// A browser strips leading and trailing controls and spaces before it parses, so strip
	// them first or every rule below judges a different string than the browser will.
	target = strings.TrimFunc(target, isControlOrSpace)

	// Redirect/link targets are often built from remote-influenced or caller-supplied
	// data, so a "javascript:" or "data:" scheme would turn an open redirect into XSS.
	parsed, err := url.Parse(target)

	// Fail closed on unparseable input, which is also what rejects an interior tab or
	// newline that a browser would strip and then re-read.
	if err != nil {
		return nil, false
	}

	// An absolute URL is safe only with an http(s) scheme; the off-site host is allowed.
	if parsed.Scheme != "" {
		return parsed, IsSchemeValid(parsed.Scheme)
	}

	// RULE: A schemeless value must be a true relative path. Two leading slashes in ANY
	// mix of "/" and "\" open an authority in a browser, whatever net/url makes of them,
	// so the target is off-site. See AGENTS.md.
	if hasAuthorityPrefix(target) {
		return parsed, false
	}

	// Belt and braces: net/url reports a host for the "//host" spelling on its own.
	return parsed, (parsed.Host == "")
}

// hasAuthorityPrefix returns TRUE if target opens with the two characters that send a
// browser into authority (hostname) parsing.
func hasAuthorityPrefix(target string) bool {
	return (len(target) >= 2) && isSlash(target[0]) && isSlash(target[1])
}

// isSlash returns TRUE if b is a character that a browser reads as a URL path separator
func isSlash(b byte) bool {
	return strings.IndexByte(authoritySlashes, b) >= 0
}

// isControlOrSpace returns TRUE if r is an ASCII C0 control character or a space
func isControlOrSpace(r rune) bool {
	return r <= ' '
}
