package dns

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHost_PortedFromDomainPackage(t *testing.T) {

	do := func(input string, expected string) {
		require.Equal(t, expected, Host(input))
	}

	do("http://localhost", "http://localhost")
	do("http://localhost/", "http://localhost")
	do("http://localhost:8080/some-path-name", "http://localhost:8080")
	do("http://localhost/some-path-name", "http://localhost")

	do("https://server.com", "https://server.com")
	do("https://server.com/", "https://server.com")
	do("https://server.com/some-path-name", "https://server.com")
	do("https://server.com/many/path/names", "https://server.com")
	do("https://server.com:443", "https://server.com:443")

	do("HTTPS://CaSeInSeNsItIvE.com", "HTTPS://CaSeInSeNsItIvE.com")

	do("not-a-host", "")
}
