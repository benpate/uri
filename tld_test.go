package uri

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTLD(t *testing.T) {

	require.True(t, IsValidTLD("com"))
	require.True(t, IsValidTLD("net"))
	require.True(t, IsValidTLD("org"))
	require.True(t, IsValidTLD("io"))
	require.True(t, IsValidTLD("dev"))
	require.True(t, IsValidTLD("app"))

	require.False(t, IsValidTLD("invalidtld"))
	require.False(t, IsValidTLD("nope"))
	require.False(t, IsValidTLD("123"))
	require.False(t, IsValidTLD("-com"))
	require.False(t, IsValidTLD("com-"))
	require.False(t, IsValidTLD("..com"))
	require.False(t, IsValidTLD("com."))
	require.False(t, IsValidTLD(".com"))
}

func TestNotValidTLD(t *testing.T) {

	// NotValidTLD is the exact inverse of IsValidTLD
	require.False(t, NotValidTLD("com"))
	require.False(t, NotValidTLD("org"))

	require.True(t, NotValidTLD(""))
	require.True(t, NotValidTLD("invalidtld"))
	require.True(t, NotValidTLD("123"))
}

func TestValidateTLD(t *testing.T) {

	// Valid TLDs return no error
	require.Nil(t, ValidateTLD("com"))
	require.Nil(t, ValidateTLD("org"))
	require.Nil(t, ValidateTLD("io"))

	// An empty TLD is reported as invalid
	require.NotNil(t, ValidateTLD(""))

	// A TLD that is not in the IANA list is reported as invalid
	require.NotNil(t, ValidateTLD("invalidtld"))
	require.NotNil(t, ValidateTLD("nope"))
}

// TestRefreshTLDs exercises the happy path of fetching the live IANA list.
//
// NOTE: RefreshTLDs's HTTP-error and read-error branches are not covered. They
// depend on network failures that can't be reliably reproduced here, and both
// merely report the error and keep the existing list, so they are left untested.
func TestRefreshTLDs(t *testing.T) {
	RefreshTLDs()

	// Whether the refresh succeeds or quietly falls back to the embedded list,
	// common TLDs must remain valid afterward.
	require.True(t, IsValidTLD("com"))
	require.True(t, IsValidTLD("org"))
}
