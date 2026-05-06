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
