package dns

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidIPAddress(t *testing.T) {
	require.True(t, IsValidIPAddress("127.0.0.1"))
	require.True(t, IsValidIPAddress("::1"))
	require.False(t, IsValidIPAddress("256.256.256.256"))
	require.False(t, IsValidIPAddress("12345::6789"))
	require.False(t, IsValidIPAddress("hamburger"))
}
