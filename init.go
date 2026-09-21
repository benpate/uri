package uri

import (
	"bytes"
	"embed"
	"regexp"
	"strings"
	"sync/atomic"
)

// embeddedFiles holds the IANA TLD list that is compiled into the binary.
//
//go:embed all:_iana.txt
var embeddedFiles embed.FS

// validTLDs holds the current set of valid IANA TLDs, which RefreshTLDs replaces
// with a single atomic swap while other goroutines may be reading it.
var validTLDs atomic.Pointer[map[string]struct{}]

// validDomainSegment matches a single RFC 1035 domain label: alphanumeric at
// both ends, with hyphens allowed in between, and never longer than 63 characters.
var validDomainSegment *regexp.Regexp = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

// init loads the embedded IANA TLD list into memory before any caller can use the package
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
func parseTLDs(data []byte) map[string]struct{} {

	// Kept separate from importTLDs so that RefreshTLDs can inspect the result
	// before it replaces a known-good list with it.

	// Split the file into lines and prepare a result map.
	lines := bytes.Split(data, []byte("\n"))
	result := make(map[string]struct{}, len(lines))

	// Keep only the lines that are a bare, valid TLD
	for _, line := range lines {

		// Normalize the TLD
		tld := strings.ToLower(string(line))

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
