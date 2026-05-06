package uri

// IsLoopback returns TRUE if the provided domain is a loopback address (localhost, 127.0.0.1, ::1)
func IsLoopback(domain string) bool {
	return (domain == "localhost") || (domain == "127.0.0.1") || (domain == "::1")
}

// NotLoopback returns TRUE if the provided domain is NOT a loopback address (otherserver.com, 192.168.0.5)
func NotLoopback(domain string) bool {
	return !IsLoopback(domain)
}
