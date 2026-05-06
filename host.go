package uri

import (
	"regexp"
)

var hostRegexp = regexp.MustCompile("(?i)^https?://[^/]+")

// Host returns the protocol and hostname
func Host(url string) string {

	if result := hostRegexp.FindStringSubmatch(url); len(result) > 0 {
		return result[0]
	}

	return ""
}
