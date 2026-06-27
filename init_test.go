package uri

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// restoreTLDs snapshots the global TLD map and registers a cleanup that
// restores it, so tests that mutate the map via importTLDs don't leak state
// into other tests.
func restoreTLDs(t *testing.T) {
	t.Helper()
	original := validTLDs.Load()
	t.Cleanup(func() {
		validTLDs.Store(original)
	})
}

// NOTE: init()'s panic branch (a failed read of the embedded _iana.txt) is not
// covered. The file is embedded at compile time, so the read cannot fail at
// runtime; the branch is unreachable and left untested.

func TestImportTLDs(t *testing.T) {

	restoreTLDs(t)

	// A typical IANA file: a comment line, then one TLD per line.
	importTLDs([]byte("# Version 2026, Last Updated ...\nCOM\nORG\nNET\n"))

	// Comments are skipped; TLDs are lower-cased and present.
	require.True(t, IsValidTLD("com"))
	require.True(t, IsValidTLD("org"))
	require.True(t, IsValidTLD("net"))

	// The comment line is not imported as a TLD.
	require.False(t, IsValidTLD("# version 2026, last updated ..."))

	// TLDs from the previous (real) list are gone — the map is fully replaced.
	require.False(t, IsValidTLD("io"))
}

func TestImportTLDs_SkipsInvalidLines(t *testing.T) {

	restoreTLDs(t)

	importTLDs([]byte("good\n\n# comment\n-bad\nbad-\nalso.bad\nUPPER\n"))

	// Valid single-segment labels are kept (and lower-cased)
	require.True(t, IsValidTLD("good"))
	require.True(t, IsValidTLD("upper"))

	// Empty lines, comments, and labels failing the segment regex are skipped
	require.False(t, IsValidTLD(""))
	require.False(t, IsValidTLD("# comment"))
	require.False(t, IsValidTLD("-bad"))
	require.False(t, IsValidTLD("bad-"))
	require.False(t, IsValidTLD("also.bad")) // contains a dot, not a single segment
}

func TestImportTLDs_Empty(t *testing.T) {

	restoreTLDs(t)

	// Importing an empty payload yields no valid TLDs (but must not panic)
	importTLDs([]byte(""))
	require.False(t, IsValidTLD("com"))
}

// TestImportTLDs_ConcurrentReadWrite republishes the TLD map (the same atomic
// swap RefreshTLDs performs) while many goroutines read it. Run under -race,
// this guards against the data race that a plain map assignment would cause.
func TestImportTLDs_ConcurrentReadWrite(t *testing.T) {

	restoreTLDs(t)

	const writers, readers, iterations = 4, 16, 500

	var wg sync.WaitGroup

	// Writers continuously republish the map via importTLDs.
	for range writers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range iterations {
				importTLDs([]byte("COM\nORG\nNET\n"))
			}
		}()
	}

	// Readers continuously validate TLDs while the map is being swapped.
	for range readers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range iterations {
				_ = IsValidTLD("com")
				_ = ValidateTLD("org")
			}
		}()
	}

	wg.Wait()
}

// FuzzImportTLDs verifies that the parser never panics on arbitrary input and
// that every TLD it accepts is a valid, lower-cased single-segment label.
func FuzzImportTLDs(f *testing.F) {

	f.Add("com\norg\n")
	f.Add("# comment\nNET\n")
	f.Add("")
	f.Add("\n\n\n")
	f.Add("-bad\nbad-\na.b\n")

	original := validTLDs.Load()
	f.Cleanup(func() {
		validTLDs.Store(original)
	})

	f.Fuzz(func(t *testing.T, data string) {

		importTLDs([]byte(data))

		for tld := range *validTLDs.Load() {
			require.True(t, validDomainSegment.MatchString(tld), "invalid TLD accepted: %q", tld)
			require.Equal(t, tld, strings.ToLower(tld), "TLD was not lower-cased: %q", tld)
		}
	})
}
