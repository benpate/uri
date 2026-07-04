package uri

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {
	require.Equal(t, "/", Path("https://localhost"))
	require.Equal(t, "/", Path("https://localhost/"))
	require.Equal(t, "/path", Path("https://localhost/path"))
	require.Equal(t, "/path", Path("https://localhost/path?and=query"))
}

func TestPath_EdgeCases(t *testing.T) {

	// A path on a non-http(s) scheme is still extracted
	require.Equal(t, "/a/b", Path("ftp://server.com/a/b"))

	// An unparseable URL (bad percent-escape) falls back to the root path
	require.Equal(t, "/", Path("https://localhost/%zz"))
}

func TestPathAndQuery(t *testing.T) {
	require.Equal(t, "/", PathAndQuery("http://localhost"))
	require.Equal(t, "/", PathAndQuery("http://localhost/"))
	require.Equal(t, "/path", PathAndQuery("http://localhost/path"))
	require.Equal(t, "/path?and=query", PathAndQuery("http://localhost/path?and=query"))
}

func TestPathAndQuery_EdgeCases(t *testing.T) {

	// Query is preserved even when the path is empty (defaults to "/")
	require.Equal(t, "/?and=query", PathAndQuery("http://localhost?and=query"))

	// An unparseable URL (bad percent-escape) falls back to the root path
	require.Equal(t, "/", PathAndQuery("https://localhost/%zz"))
}

// FuzzPath verifies that Path and PathAndQuery never panic on arbitrary input.
func FuzzPath(f *testing.F) {

	f.Add("https://localhost/path?and=query")
	f.Add("http://localhost")
	f.Add("not a url")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		path := Path(input)
		pathAndQuery := PathAndQuery(input)

		// Both functions always return a non-empty, root-absolute path.
		require.True(t, strings.HasPrefix(path, "/"),
			"Path(%q) = %q must be root-absolute", input, path)
		require.True(t, strings.HasPrefix(pathAndQuery, "/"),
			"PathAndQuery(%q) = %q must be root-absolute", input, pathAndQuery)

		// PathAndQuery always begins with whatever Path returned.
		require.True(t, strings.HasPrefix(pathAndQuery, path),
			"PathAndQuery %q should start with Path %q", pathAndQuery, path)
	})
}
