package uri

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// isProtocolRelative reports whether a browser would treat s as pointing at another host
// when placed in a Location header.  It mirrors the browser rule that "\" is equivalent
// to "/" in the authority position, so it checks the slash-normalized form.
func isProtocolRelative(s string) bool {
	normalized := strings.ReplaceAll(s, "\\", "/")
	return strings.HasPrefix(normalized, "//")
}

// TestPathAndQuery_Redirect covers PathAndQuery as an open-redirect sanitizer: every
// off-site vector must be neutralized to a same-origin path, and every legitimate local
// path must be preserved exactly.
func TestPathAndQuery_Redirect(t *testing.T) {

	// safe pins the exact, expected output for one input.
	safe := func(name, input, want string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			require.Equal(t, want, PathAndQuery(input), "PathAndQuery(%q)", input)
		})
	}

	// neutralized asserts an input is rendered non-off-site: the result is a root-absolute
	// path that is not protocol-relative and carries no authority.  Used for attack vectors
	// where the exact same-origin path is unimportant, only that it can no longer redirect
	// to another host.
	neutralized := func(name, input string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			result := PathAndQuery(input)

			require.True(t, strings.HasPrefix(result, "/"),
				"PathAndQuery(%q) = %q must be root-absolute", input, result)
			require.False(t, isProtocolRelative(result),
				"PathAndQuery(%q) = %q is an open-redirect vector", input, result)

			// The neutralized result must re-parse to a URL with no authority.
			parsed, err := url.Parse(result)
			require.NoError(t, err, "PathAndQuery(%q) = %q should itself be parseable", input, result)
			require.Empty(t, parsed.Host,
				"PathAndQuery(%q) = %q re-parses to host %q", input, result, parsed.Host)
		})
	}

	// --- Open-redirect vectors: must be neutralized ---
	neutralized("authority-relative double slash", "//evil.com")
	neutralized("triple slash", "///evil.com")
	neutralized("many slashes", "////evil.com")
	neutralized("percent-encoded double slash", "/%2f%2fevil.com")
	neutralized("percent-encoded single slash", "/%2fevil.com")
	neutralized("backslash authority", `/\evil.com`)
	neutralized("mixed slash backslash", `/\/evil.com`)
	neutralized("double backslash", `/\\evil.com`)
	neutralized("slash then backslash", `//\evil.com`)
	neutralized("leading double backslash", `\\evil.com`)
	neutralized("backslash then slash", `\/evil.com`)
	neutralized("percent-encoded backslashes", "/%5c%5cevil.com")
	neutralized("authority with path and query", "//evil.com/path?q=1")
	neutralized("absolute url strips host", "https://evil.com/x")
	neutralized("uppercase scheme strips host", "HTTPS://evil.com/x")

	// --- Legitimate local references: preserved exactly ---
	safe("simple path", "/dashboard", "/dashboard")
	safe("at-username path", "/@me/settings", "/@me/settings")
	safe("path with query", "/a/b?x=1", "/a/b?x=1")
	safe("path and query", "/path?and=query", "/path?and=query")
	safe("root", "/", "/")
	safe("empty defaults to root", "", "/")
	safe("query only", "?justquery=1", "/?justquery=1")
	safe("trailing slash", "/trailing/", "/trailing/")
	safe("absolute local url", "http://localhost/path", "/path")
	safe("absolute local url with query", "https://localhost/a/b?x=1", "/a/b?x=1")

	// A double slash or backslash INSIDE the path (not at the front) is legitimate.
	safe("internal double slash preserved", "/double//slash/inside", "/double//slash/inside")
	safe("internal backslash preserved", "/dir/sub\\file", "/dir/sub\\file")

	// A path with no leading slash is forced to be root-absolute.
	safe("no leading slash forced absolute", "path/no/leading/slash", "/path/no/leading/slash")
	safe("bare word forced absolute", "relative", "/relative")
	safe("bare word with query forced absolute", "relative?q=1", "/relative?q=1")

	// An unparseable URL falls back to the root path (never empty).
	safe("bad percent-escape falls back to root", "https://localhost/%zz", "/")
	safe("bad host falls back to root", "http://[::1", "/")
}

// TestPath_Redirect mirrors TestPathAndQuery_Redirect for Path, which shares the same
// leading-authority hardening (and therefore the same open-redirect guarantees).
func TestPath_Redirect(t *testing.T) {

	safe := func(name, input, want string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			require.Equal(t, want, Path(input), "Path(%q)", input)
		})
	}

	neutralized := func(name, input string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			result := Path(input)

			require.True(t, strings.HasPrefix(result, "/"),
				"Path(%q) = %q must be root-absolute", input, result)
			require.False(t, isProtocolRelative(result),
				"Path(%q) = %q is an open-redirect vector", input, result)
		})
	}

	neutralized("triple slash", "///evil.com")
	neutralized("backslash authority", `/\evil.com`)
	neutralized("leading double backslash", `\\evil.com`)
	neutralized("authority-relative double slash", "//evil.com")

	safe("simple path", "/dashboard", "/dashboard")
	safe("strips host and query", "https://localhost/path?and=query", "/path")
	safe("no leading slash forced absolute", "path/no/slash", "/path/no/slash")
	safe("empty defaults to root", "", "/")
	safe("bad percent-escape falls back to root", "https://localhost/%zz", "/")
}

// FuzzPathAndQuery_SafeRedirect is the security invariant under fuzzing: for ANY input,
// the result of PathAndQuery is never something a browser would treat as an off-site
// redirect.  It is always a root-absolute path that is not protocol-relative and carries
// no authority.
func FuzzPathAndQuery_SafeRedirect(f *testing.F) {

	seeds := []string{
		"", "/", "//evil.com", "///evil.com", "/%2f%2fevil.com", `/\evil.com`,
		`\\evil.com`, "https://evil.com/x", "/legit/path?q=1", "path/no/slash",
		"?q=1", "/a//b", "javascript:alert(1)", "//evil.com/x?y=z#frag",
		"/%5c%5cevil.com", "\t//evil.com",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {

		result := PathAndQuery(input)

		// INVARIANT 1: the result is always a non-empty, root-absolute path.
		require.True(t, strings.HasPrefix(result, "/"),
			"PathAndQuery(%q) = %q is not rooted at '/'", input, result)

		// INVARIANT 2: the result is never protocol-relative (no off-site redirect),
		// including under the browser's backslash-as-slash interpretation.
		require.False(t, isProtocolRelative(result),
			"PathAndQuery(%q) = %q is protocol-relative (open redirect)", input, result)

		// INVARIANT 3: whatever we return must itself re-parse to a URL with NO host,
		// so it can never carry an authority into a same-origin request.
		if parsed, err := url.Parse(result); err == nil {
			require.Empty(t, parsed.Host,
				"PathAndQuery(%q) = %q re-parses to a non-empty host %q", input, result, parsed.Host)
		}
	})
}
