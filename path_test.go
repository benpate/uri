package dns

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {
	require.Equal(t, "/", Path("https://localhost"))
	require.Equal(t, "/", Path("https://localhost/"))
	require.Equal(t, "/path", Path("https://localhost/path"))
	require.Equal(t, "/path", Path("https://localhost/path?and=query"))
}

func TestPathAndQuery(t *testing.T) {
	require.Equal(t, "/", PathAndQuery("http://localhost"))
	require.Equal(t, "/", PathAndQuery("http://localhost/"))
	require.Equal(t, "/path", PathAndQuery("http://localhost/path"))
	require.Equal(t, "/path?and=query", PathAndQuery("http://localhost/path?and=query"))
}
