package uri

import (
	"net"
	"net/url"
	"strings"
)

// NormalizeHost extracts the bare hostname from any URL-ish value, regardless
// of its shape: a bare hostname, a "host:port", a full "scheme://" URL, a value
// with userinfo ("user@host"), or a bracketed IPv6 literal ("[::1]:8080"). The
// result is lower-cased. Unlike Hostname, the result MAY contain colons, because
// an IPv6 address (e.g. "::1") is returned in its canonical, unbracketed form.
func NormalizeHost(value string) string {

	value = strings.TrimSpace(value)

	if value == "" {
		return ""
	}

	// If the value parses as a URL *with a scheme*, trust net/url to pull out the
	// host. This correctly handles userinfo and bracketed IPv6 literals.
	if parsed, err := url.Parse(value); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		return strings.ToLower(parsed.Hostname())
	}

	// No scheme: strip any path, query, or fragment that follows the authority.
	if index := strings.IndexAny(value, "/?#"); index >= 0 {
		value = value[:index]
	}

	// Strip userinfo ("user:pass@host" -> "host"). The host is whatever follows
	// the last "@", so credentials containing "@" are handled correctly.
	if at := strings.LastIndex(value, "@"); at >= 0 {
		value = value[at+1:]
	}

	value = strings.ToLower(value)

	// Bracketed IPv6, with or without a port: "[::1]:8080" -> "::1"
	if strings.HasPrefix(value, "[") {
		if end := strings.Index(value, "]"); end >= 0 {
			return value[1:end]
		}
	}

	// A single colon means "host:port" — strip the port. Multiple colons mean an
	// unbracketed IPv6 literal (e.g. "::1"), which is returned unchanged.
	if strings.Count(value, ":") == 1 {
		if host, _, err := net.SplitHostPort(value); err == nil {
			return host
		}
		value, _, _ = strings.Cut(value, ":")
	}

	return value
}
