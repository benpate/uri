package uri

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTLD verifies that IsValidTLD accepts listed TLDs and rejects unknown or malformed labels
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

// TestNotValidTLD verifies that NotValidTLD is the inverse of IsValidTLD
func TestNotValidTLD(t *testing.T) {

	require.False(t, NotValidTLD("com"))
	require.False(t, NotValidTLD("org"))

	require.True(t, NotValidTLD(""))
	require.True(t, NotValidTLD("invalidtld"))
	require.True(t, NotValidTLD("123"))
}

// TestValidateTLD verifies that ValidateTLD returns an error for an empty or unlisted TLD
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

// TestRefreshTLDs verifies that common TLDs stay valid after a refresh from the live IANA list
func TestRefreshTLDs(t *testing.T) {
	RefreshTLDs()

	// Whether the refresh succeeds or quietly falls back to the embedded list,
	// common TLDs must remain valid afterward.
	require.True(t, IsValidTLD("com"))
	require.True(t, IsValidTLD("org"))
}

// TestRefreshTLDs_KeepsListOnHTTPError verifies that an error status from IANA leaves the TLD list intact
func TestRefreshTLDs_KeepsListOnHTTPError(t *testing.T) {

	restoreTLDs(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("<html><body>503 Service Unavailable</body></html>"))
	}))
	defer server.Close()

	refreshTLDsFrom(server.URL)

	require.True(t, IsValidTLD("com"))
	require.True(t, IsValidTLD("org"))
	require.Nil(t, ValidateURL("https://example.com"))
}

// TestRefreshTLDs_KeepsListOnEmptyBody verifies that a success response with an empty body
// leaves the TLD list intact
func TestRefreshTLDs_KeepsListOnEmptyBody(t *testing.T) {

	restoreTLDs(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	refreshTLDsFrom(server.URL)

	require.True(t, IsValidTLD("com"))
	require.Nil(t, ValidateURL("https://example.com"))
}

// TestRefreshTLDs_KeepsListOnNonTLDBody verifies that a success response with no valid TLDs,
// such as a captive-portal page, leaves the TLD list intact
func TestRefreshTLDs_KeepsListOnNonTLDBody(t *testing.T) {

	restoreTLDs(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("<!DOCTYPE html>\n<html>\n<head><title>Sign in</title></head>\n</html>\n"))
	}))
	defer server.Close()

	refreshTLDsFrom(server.URL)

	require.True(t, IsValidTLD("com"))
	require.Nil(t, ValidateURL("https://example.com"))
}

// TestRefreshTLDs_PublishesGoodList verifies that a well-formed IANA response replaces the TLD list
func TestRefreshTLDs_PublishesGoodList(t *testing.T) {

	restoreTLDs(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("# Version 2026, Last Updated ...\nCOM\nORG\nEXAMPLETLD\n"))
	}))
	defer server.Close()

	refreshTLDsFrom(server.URL)

	// The new list is live
	require.True(t, IsValidTLD("exampletld"))

	// ...and it fully replaced the old one
	require.False(t, IsValidTLD("io"))
}

// TestRefreshTLDs_KeepsListOnConnectionFailure verifies that a transport error leaves the TLD list intact
func TestRefreshTLDs_KeepsListOnConnectionFailure(t *testing.T) {

	restoreTLDs(t)

	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	serverURL := server.URL
	server.Close()

	refreshTLDsFrom(serverURL)

	require.True(t, IsValidTLD("com"))
}
