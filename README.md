# uri

Resilient helpers for parsing, validating, and inspecting URLs and hostnames: scheme/protocol handling, IANA top-level-domain validation, local-network detection, and SSRF-safe public-IP classification.

Most functions are forgiving by design — they take a raw string, do their best, and return an empty string or `false` rather than an error when the input is malformed. The `Validate*` functions are the exception: they return a `derp` error describing exactly what was wrong.

## Map for the impatient

- **URLs** — [url.go](url.go): `ParseURL` (a stricter `url.Parse` that *requires* an http/https scheme and a valid hostname), `ValidateURL`, `IsValidURL`.
- **Hostnames** — [hostname.go](hostname.go): `Hostname` (strip everything down to the host), `ValidateHostname` (RFC 1035 lengths + IANA TLD check, with IP and local names exempted).
- **Scheme / protocol** — [scheme.go](scheme.go), [protocol.go](protocol.go): `Scheme` (no `://`) vs `Protocol` (with `://`); `GuessSchemeForHostname` / `PrependProtocol` pick http for local hosts and https for public ones.
- **Local-network detection** — [localhost.go](localhost.go), [loopback.go](loopback.go): `IsLocalHostname`, `IsLocalURL`, `IsLoopback`.
- **Public-IP / SSRF** — [public.go](public.go): `IsPublicIP`, `NotPublicIP`, `IsPublicIPAddress`.
- **TLDs** — [tld.go](tld.go), [init.go](init.go): `IsValidTLD`, `ValidateTLD`, `RefreshTLDs`; the IANA list is embedded from [_iana.txt](_iana.txt).

## What matters here

- **`IsPublicIP` (public.go) and `IsLocalHostname` (localhost.go) are two different trust boundaries — don't substitute one for the other.** `IsLocalHostname` answers "should this use http and skip TLS" (loopback, `.local`, RFC 1918, link-local, unspecified). `IsPublicIP` answers "is this safe to *connect to*" and is the SSRF gate — it rejects everything `IsLocalHostname` does **plus** the full IANA special-use registry (CGNAT, TEST-NETs, 6to4/NAT64 ranges that can embed a private IPv4, etc.) that the stdlib predicates miss. For an SSRF check, always use `IsPublicIP` on the *resolved* IP at connection time (e.g. in a `net.Dialer.Control` hook), never a hostname string check.

- **The cloud-metadata endpoint `169.254.169.254` is deliberately caught by the link-local case** in both `IsLocalHostname` and `IsPublicIP`. It looks like an ordinary public-ish address but routes to instance credentials — the single most important SSRF target. Don't "simplify" the link-local branches away.

- **`IsPublicIP` works on the resolved `net.IP`, not on a string — and that's the point.** A string like `0x7f.1` or `2130706433` is loopback once parsed, but a naive string check won't see it. Resolve first, classify second.

- **`ValidateHostname` exempts IP addresses and local names *before* applying DNS rules.** `127.0.0.1` and `friday.local` are valid hostnames here even though they have no IANA TLD. Only after those exemptions does it enforce RFC 1035 label/length limits and require the final segment to be a real IANA TLD. If you change the ordering, you'll start rejecting valid loopback/local inputs.

- **The TLD list is embedded at build time and loaded once in `init()`.** `RefreshTLDs` can re-fetch the live list from IANA at runtime, but it is best-effort: on any network/read error it logs via `derp` and silently keeps the embedded list. It also caps the download with `io.LimitReader(…, 1<<20)` — keep that cap if you touch it.

- **Paired `Is…` / `Not…` predicates are intentional, not redundant.** `NotLocalURL`, `NotPublicIP`, `NotValidTLD`, etc. exist so callers read naturally at the call site (`if uri.NotPublicIP(ip)`). Each `Not…` is a one-line negation of its `Is…` twin — keep them in sync.

- **`Hostname` lower-cases and strips, it does not validate.** It will happily return garbage from garbage. Run the result through `ValidateHostname` if the input is untrusted.

## See also

- [benpate/derp](../derp/README.md) — the error type returned by the `Validate*` functions.
