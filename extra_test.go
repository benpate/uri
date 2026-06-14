package uri

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateHostname_TooLong(t *testing.T) {

	// A domain longer than 253 characters is invalid (RFC 1035).
	// Build it from valid 50-char segments so length is the only failure.
	segment := strings.Repeat("a", 50)
	hostname := strings.Join([]string{segment, segment, segment, segment, segment, "com"}, ".")
	require.Greater(t, len(hostname), 253)

	err := ValidateHostname(hostname)
	require.Error(t, err)
}

func TestValidateHostname_SegmentTooLong(t *testing.T) {

	// A single label longer than 63 characters is invalid (RFC 1035)
	hostname := strings.Repeat("a", 64) + ".com"
	err := ValidateHostname(hostname)
	require.Error(t, err)
}

func TestValidateHostname_Valid(t *testing.T) {
	require.NoError(t, ValidateHostname("example.com"))
	require.NoError(t, ValidateHostname("localhost"))
	require.NoError(t, ValidateHostname("127.0.0.1"))
}
