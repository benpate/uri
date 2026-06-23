package uri

import (
	"bytes"
	"embed"
	"regexp"
	"strings"
)

//go:embed all:_iana.txt
var embeddedFiles embed.FS
var validTLDs map[string]struct{}
var validDomainSegment *regexp.Regexp = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

func init() {

	// Read the embedded file for TLD information
	data, err := embeddedFiles.ReadFile("_iana.txt")
	if err != nil {
		panic(err)
	}

	// Import the TLDs into the memory cache
	importTLDs(data)
}

func importTLDs(data []byte) {

	// Split the file into lines and prepare a result map.
	lines := bytes.Split(data, []byte("\n"))
	result := make(map[string]struct{}, len(lines)-2)

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

	// Copy the TLDs into the global variable
	validTLDs = result
}
