# 🧭 uri

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/uri.svg)](https://pkg.go.dev/github.com/benpate/uri)
[![Version](https://img.shields.io/github/v/release/benpate/uri?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/uri/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/uri/go.yml?branch=main)](https://github.com/benpate/uri/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/uri?style=flat-square)](https://goreportcard.com/report/github.com/benpate/uri)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/uri.svg?style=flat-square)](https://codecov.io/gh/benpate/uri)

## Resilient URL and Hostname Helpers for Go

Resilient helpers for parsing, validating, and inspecting URLs and hostnames: scheme/protocol handling, IANA top-level-domain validation, local-network detection, and SSRF-safe public-IP classification.

Most functions are forgiving by design — they take a raw string, do their best, and return an empty string or `false` rather than an error when the input is malformed. The `Validate*` functions are the exception: they return a `derp` error describing exactly what was wrong.

## Map for the Impatient

- **URLs** — [url.go](url.go): `ParseURL` (a stricter `url.Parse` that *requires* an http/https scheme and a valid hostname), `ValidateURL`, `IsValidURL`.
- **Hostnames** — [hostname.go](hostname.go): `Hostname` (strip everything down to the host), `ValidateHostname` (RFC 1035 lengths + IANA TLD check, with IP and local names exempted).
- **Scheme / protocol** — [scheme.go](scheme.go), [protocol.go](protocol.go): `Scheme` (no `://`) vs `Protocol` (with `://`, and `""` when the value never had one); `GuessSchemeForHostname` / `PrependProtocol` pick http for local hosts and https for public ones.
- **Origins** — [host.go](host.go), [origin.go](origin.go): `Host` returns the RFC 6454 origin (default ports omitted); `IsSameOrigin` / `NotSameOrigin` compare two of them.
- **Redirect / link safety** — [redirect.go](redirect.go), [safeurl.go](safeurl.go): `IsSafeRedirectURL` / `NotSafeRedirectURL` accept only same-site relative paths or absolute http(s) URLs, rejecting `javascript:`/`data:`/protocol-relative targets that turn an open redirect into XSS. Targets are judged as a **browser** parses them, so the spellings `net/url` reads as same-site and a browser reads as off-site (`/\evil.com`, `\\evil.com`, `///evil.com`, `" //evil.com"`) are rejected too. `IsValidRedirectURL` / `NotValidRedirectURL` add *absolute, with a valid host*, which is the rule for a destination named by a remote server. `SafeURL` applies the base rule and additionally *neutralizes* the value for use in an HTML attribute.
- **Same-origin paths** — [path.go](path.go): `Path` and `PathAndQuery` reduce any URL to a root-absolute path on this host, for a `Location` header or a post-login `next`. Like the redirect checks, they neutralize what a **browser** would read as off-site, including an encoded tab or newline (`/%09/evil.com`).
- **Local-network detection** — [localhost.go](localhost.go), [loopback.go](loopback.go): `IsLocalHostname`, `IsLocalURL`, `IsLoopback`.
- **Public-IP / SSRF** — [public.go](public.go): `IsPublicIP`, `NotPublicIP`, `IsPublicIPAddress`. Apply `IsPublicIP` to the *resolved* IP at connection time (e.g. in a `net.Dialer` `Control` hook), never to a hostname string.
- **TLDs** — [tld.go](tld.go), [init.go](init.go): `IsValidTLD`, `ValidateTLD`, `RefreshTLDs`; the IANA list is embedded from [_iana.txt](_iana.txt).

## See Also

- [benpate/derp](../derp/README.md) — the error type returned by the `Validate*` functions.

## Pull Requests Welcome

I'm trying to make uri the best it can be, and your help is greatly appreciated. If you find a bug or have an idea for a new feature, please open an issue or submit a pull request. We're all in this together! 🧭
