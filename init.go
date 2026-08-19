package uri

import (
	"bytes"
	"embed"
	"regexp"
	"strings"
	"sync/atomic"
)

// embeddedFiles holds the IANA TLD list that is compiled into the binary,
// so that the package is fully functional with no network access.
//
//go:embed all:_iana.txt
var embeddedFiles embed.FS

// validTLDs holds the set of valid IANA TLDs. It is an atomic pointer because
// RefreshTLDs may replace the whole map at runtime while other goroutines are
// reading it; an atomic swap keeps reads lock-free and free of data races.
var validTLDs atomic.Pointer[map[string]struct{}]

// validDomainSegment matches a single RFC 1035 domain label: alphanumeric at
// both ends, with hyphens allowed in between, and never longer than 63 characters.
var validDomainSegment *regexp.Regexp = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

// init loads the embedded IANA TLD list into memory before any caller can use
// the package, so that TLD validation never depends on a network round-trip.
func init() {

	// Read the embedded file for TLD information
	data, err := embeddedFiles.ReadFile("_iana.txt")
	if err != nil {
		panic(err)
	}

	// Import the TLDs into the memory cache
	importTLDs(data)
}

// importTLDs parses an IANA TLD list and publishes it as the in-memory set,
// replacing whatever was there before.
func importTLDs(data []byte) {

	result := parseTLDs(data)

	// Publish the new map with a single atomic swap so concurrent readers see
	// either the old map or the new one, never a partially-built map.
	validTLDs.Store(&result)
}

// parseTLDs reads an IANA TLD list into a set of lower-cased top-level domains,
// skipping comments, blank lines, and anything that is not a single valid label.
// Parsing is kept separate from publishing so that a caller can inspect the
// result -- see RefreshTLDs -- before replacing a known-good list with it.
func parseTLDs(data []byte) map[string]struct{} {

	// Split the file into lines and prepare a result map.
	lines := bytes.Split(data, []byte("\n"))
	result := make(map[string]struct{}, len(lines))

	for _, line := range lines {

		// Normalize the TLD
		tld := strings.ToLower((string(line)))

		// Skip if this is not a valid domain segment
		// (for instance, a comment or empty line)
		if !validDomainSegment.MatchString(tld) {
			continue
		}

		// Add the TLD to the map of valid TLDs
		result[tld] = struct{}{}
	}

	return result
}
