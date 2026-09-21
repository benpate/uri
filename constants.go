package uri

// SchemeHTTP represents the "http" network scheme (without the :// suffix)
const SchemeHTTP = "http"

// SchemeHTTPS represents the "https" network scheme (without the :// suffix)
const SchemeHTTPS = "https"

// ProtocolHTTP represents the "http://" protocol string (including the :// suffix)
const ProtocolHTTP = "http://"

// ProtocolHTTPS represents the "https://" protocol string (including the :// suffix)
const ProtocolHTTPS = "https://"

// ProtocolSuffix represents the "://" suffix that follows the protocol name
const ProtocolSuffix = "://"

// ianaTLDList is the canonical location of the IANA top-level-domain list.
const ianaTLDList = "https://data.iana.org/TLD/tlds-alpha-by-domain.txt"

// authoritySlashes are the characters that open a URL authority, which a browser reads
// alike for http(s).
const authoritySlashes = `/\`

// tabOrNewline are the characters that a browser deletes from anywhere in a URL before
// it parses one.
const tabOrNewline = "\t\n\r"
