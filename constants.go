package uri

// authoritySlashes are the characters that open a URL authority, which a browser reads
// alike for http(s). safeLocalPath collapses exactly this set, so the two cannot drift.
const authoritySlashes = `/\`
