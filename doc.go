// Package uri provides helpers for parsing, validating, and inspecting URLs
// and hostnames.
//
// Most functions are forgiving: they take a raw string of any shape, do their
// best, and return "" or FALSE rather than an error. The Validate* functions
// are the exception, returning a derp error that describes what was wrong.
// Predicates come in Is…/Not… pairs so that call sites read naturally.
//
// The package covers five areas:
//
//   - URLs and hostnames — ParseURL is a stricter url.Parse that requires an
//     http(s) scheme and a valid hostname; ValidateHostname enforces RFC 1035
//     label rules and checks the final segment against the IANA TLD list.
//   - Schemes, protocols, and origins — Scheme and Protocol pull a URL apart,
//     Host returns an RFC 6454 origin, and IsSameOrigin compares two of them.
//   - Local-network detection — IsLoopback and IsLocalHostname decide whether
//     an address is on this machine or this network, which is what chooses
//     http over https for a local host.
//   - SSRF defense — IsPublicIP classifies a resolved net.IP against the IANA
//     special-use registries. It is the gate to apply before connecting to a
//     user-supplied host, and it is deliberately stricter than IsLocalHostname.
//   - Redirect and link safety — IsSafeRedirectURL rejects the "javascript:",
//     "data:", and protocol-relative targets that turn an open redirect into
//     XSS, and SafeURL additionally escapes a value for use in an attribute.
//
// The IANA TLD list is embedded at build time, so the package needs no network
// access. RefreshTLDs can update it at runtime, and fails closed onto the list
// already in memory.
package uri
