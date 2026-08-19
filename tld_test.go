package uri

import (
	"net/http"
	"net/http/httptest"
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

// TestRefreshTLDs_KeepsListOnHTTPError proves that an IANA outage cannot wipe
// the TLD list. Without the status-code check, the body of the error page is
// parsed as the new list and every public hostname stops validating.
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

// TestRefreshTLDs_KeepsListOnEmptyBody proves that a success response with no
// usable TLDs (an empty body, or a truncated transfer) is rejected rather than
// published as an empty list.
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

// TestRefreshTLDs_KeepsListOnNonTLDBody proves that a success response carrying
// something other than a TLD list -- a captive portal or proxy interstitial --
// is rejected, because it parses to zero valid TLDs.
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

// TestRefreshTLDs_PublishesGoodList confirms that a well-formed IANA response
// still replaces the list, so the guards above do not block the happy path.
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

// TestRefreshTLDs_KeepsListOnConnectionFailure confirms the pre-existing
// transport-error path still keeps the embedded list.
func TestRefreshTLDs_KeepsListOnConnectionFailure(t *testing.T) {

	restoreTLDs(t)

	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	serverURL := server.URL
	server.Close()

	refreshTLDsFrom(serverURL)

	require.True(t, IsValidTLD("com"))
}
