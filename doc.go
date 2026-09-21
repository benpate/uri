// Package uri provides helpers for parsing, validating, and inspecting URLs and hostnames.
//
// Most functions are forgiving: they take a raw string of any shape and return "" or
// FALSE rather than an error. The Validate* functions are the exception, returning a derp
// error that describes what was wrong. Predicates come in Is…/Not… pairs.
//
//   - URLs and hostnames: ParseURL is a stricter url.Parse, and ValidateHostname enforces
//     RFC 1035 label rules and the IANA TLD list, which is embedded at build time and
//     refreshed at runtime by RefreshTLDs.
//   - Schemes, protocols, and origins: Scheme, Protocol, Host, and IsSameOrigin.
//   - Local-network detection: IsLoopback and IsLocalHostname, which choose http over
//     https for a local host.
//   - SSRF defense: IsPublicIP classifies a resolved net.IP against the IANA special-use
//     registries. Apply it before connecting to a user-supplied host.
//   - Redirect and link safety: IsSafeRedirectURL, IsValidRedirectURL, SafeURL, and
//     PathAndQuery judge a target the way a browser parses it, not the way net/url does.
package uri
