# uri — Notes for AI Agents

- **`IsPublicIP` (public.go) and `IsLocalHostname` (localhost.go) are two different trust boundaries — don't substitute one for the other.** `IsLocalHostname` answers "should this use http and skip TLS" (loopback, `.local`, RFC 1918, link-local, unspecified). `IsPublicIP` answers "is this safe to *connect to*" and is the SSRF gate — it rejects everything `IsLocalHostname` does **plus** the full IANA special-use registry (CGNAT, TEST-NETs, 6to4/NAT64 ranges that can embed a private IPv4, etc.) that the stdlib predicates miss. For an SSRF check, always use `IsPublicIP` on the *resolved* IP at connection time (e.g. in a `net.Dialer.Control` hook), never a hostname string check.

- **The cloud-metadata endpoint `169.254.169.254` is deliberately caught by the link-local case** in both `IsLocalHostname` and `IsPublicIP`. It looks like an ordinary public-ish address but routes to instance credentials — the single most important SSRF target. Don't "simplify" the link-local branches away.

- **`IsPublicIP` works on the resolved `net.IP`, not on a string — and that's the point.** A string like `0x7f.1` or `2130706433` is loopback once parsed, but a naive string check won't see it. Resolve first, classify second.

- **`ValidateHostname` exempts IP addresses and local names *before* applying DNS rules.** `127.0.0.1` and `friday.local` are valid hostnames here even though they have no IANA TLD. Only after those exemptions does it enforce RFC 1035 label/length limits and require the final segment to be a real IANA TLD. If you change the ordering, you'll start rejecting valid loopback/local inputs.

- **The TLD list is embedded at build time and loaded once in `init()`.** `RefreshTLDs` can re-fetch the live list from IANA at runtime, but it is best-effort and **fails closed onto the list already in memory**. Three guards make that true, and all three must stay: a non-200 status is rejected (an IANA 503 page would otherwise be parsed *as* the TLD list), the download is capped with `io.LimitReader(…, 1<<20)`, and a response that parses to **zero** TLDs is rejected rather than published. That last one is why parsing (`parseTLDs`) is split from publishing (`importTLDs`) — publishing an empty map makes every public hostname fail validation process-wide, so `RefreshTLDs` inspects the parse result *before* the atomic swap.

- **Paired `Is…` / `Not…` predicates are intentional, not redundant.** `NotLocalURL`, `NotPublicIP`, `NotValidTLD`, etc. exist so callers read naturally at the call site (`if uri.NotPublicIP(ip)`). Each `Not…` is a one-line negation of its `Is…` twin — keep them in sync.

- **`Hostname` lower-cases and strips, it does not validate.** It will happily return garbage from garbage. Run the result through `ValidateHostname` if the input is untrusted.

- **`Hostname` is a thin alias for `NormalizeHost` — keep it that way.** It used to have its own cut-at-the-first-delimiter logic, which left the querystring attached (`example.com?q=1`), left userinfo attached (`evil.com@good.com`), and truncated a bracketed IPv6 literal to `"["`. Callers build federated `@user@host` handles and referer checks out of this, so both functions must answer identically for every shape; `TestHostname_MatchesNormalizeHost` pins that. Note that an IPv6 result comes back **unbracketed** (`::1`), so it can contain colons.
