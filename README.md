# 🧝 uri

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
- **Scheme / protocol** — [scheme.go](scheme.go), [protocol.go](protocol.go): `Scheme` (no `://`) vs `Protocol` (with `://`); `GuessSchemeForHostname` / `PrependProtocol` pick http for local hosts and https for public ones.
- **Local-network detection** — [localhost.go](localhost.go), [loopback.go](loopback.go): `IsLocalHostname`, `IsLocalURL`, `IsLoopback`.
- **Public-IP / SSRF** — [public.go](public.go): `IsPublicIP`, `NotPublicIP`, `IsPublicIPAddress`.
- **TLDs** — [tld.go](tld.go), [init.go](init.go): `IsValidTLD`, `ValidateTLD`, `RefreshTLDs`; the IANA list is embedded from [_iana.txt](_iana.txt).

## See Also

- [benpate/derp](../derp/README.md) — the error type returned by the `Validate*` functions.

## Pull Requests Welcome

I'm trying to make uri the best it can be, and your help is greatly appreciated. If you find a bug or have an idea for a new feature, please open an issue or submit a pull request. We're all in this together! 🧝
