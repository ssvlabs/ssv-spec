package ssv

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

// BuilderEntry is one cluster-configured direct-builder connection for the builder-request-auth extension
// (SIP #94 §5). Data is the token agreed with the builder out of band, signed byte-for-byte; when empty it
// defaults to the builder URL's hostname (see AuthData). Entries MUST be byte-identical across all
// operators, and SSV caps them at MaxBuilderEntries per validator.
type BuilderEntry struct {
	Data []byte
	URL  string
}

// MaxBuilderEntries is SSV's per-validator cap on configured builder entries (SIP #94 §5) — a sub-cap of
// the beacon-API's MAX_BUILDER_ENTRIES (64) that bounds the §7 message-validation budget.
const MaxBuilderEntries = 8

// AuthData returns the bytes signed for this entry: the configured Data, or — when Data is empty — the
// default derived from URL per the upstream get_default_auth_data rule (SIP #94 §5, builder-specs #168):
// the URL's hostname, lowercased ASCII, dropping scheme, userinfo, port, path, query and fragment, with an
// IPv6 literal in RFC 5952 compressed hex-only form inside brackets. It returns nil when neither is
// available (an internationalized hostname must be supplied in punycode); the auth round skips an empty
// result as invalid.
func (e BuilderEntry) AuthData() []byte {
	if len(e.Data) > 0 {
		return e.Data
	}
	return defaultAuthData(e.URL)
}

// defaultAuthData applies the upstream get_default_auth_data rule to a builder URL, returning nil when no
// ASCII hostname can be derived (an unparsable URL or a non-ASCII host).
func defaultAuthData(rawURL string) []byte {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}
	// url.Hostname drops userinfo and the port and unbrackets an IPv6 literal.
	host := u.Hostname()
	// U+0130 (İ) is the one code point where Go's strings.ToLower and Python's str.lower disagree on the
	// result's ASCII-ness: Go folds it to ASCII "i", but upstream's Python folds it to "i" + U+0307 (a
	// combining dot) and rejects the host via isascii(). Reject it before lowercasing so the derived auth
	// root can't diverge from upstream's (SIP #94 §5, builder-specs #168).
	if strings.ContainsRune(host, '\u0130') {
		return nil
	}
	host = strings.ToLower(host)
	if host == "" || !isASCII(host) {
		return nil
	}
	if strings.Contains(host, ":") { // IPv6 literal
		ip := net.ParseIP(host)
		if ip == nil {
			return nil
		}
		host = "[" + compressIPv6(ip) + "]"
	}
	return []byte(host)
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 0x7f {
			return false
		}
	}
	return true
}

// compressIPv6 renders an IPv6 address in RFC 5952 form: lowercase hex groups with the longest run of zero
// groups collapsed to "::". Unlike net.IP.String it never uses the mixed IPv4 dotted notation, so
// ::ffff:192.0.2.1 renders as ::ffff:c000:201, matching get_default_auth_data.
func compressIPv6(ip net.IP) string {
	ip = ip.To16()
	if ip == nil {
		return ""
	}
	var groups [8]uint16
	for i := range groups {
		groups[i] = uint16(ip[2*i])<<8 | uint16(ip[2*i+1])
	}
	// Pick the longest run of zero groups (length >= 2) to collapse to "::".
	bestStart, bestLen := -1, 0
	for i := 0; i < 8; {
		if groups[i] != 0 {
			i++
			continue
		}
		j := i
		for j < 8 && groups[j] == 0 {
			j++
		}
		if j-i > bestLen {
			bestStart, bestLen = i, j-i
		}
		i = j
	}
	if bestLen < 2 {
		bestStart = -1
	}
	var b strings.Builder
	lastColon := false
	for i := 0; i < 8; i++ {
		if i == bestStart {
			b.WriteString("::")
			lastColon = true
			i += bestLen - 1
			continue
		}
		if b.Len() > 0 && !lastColon {
			b.WriteByte(':')
		}
		b.WriteString(strconv.FormatUint(uint64(groups[i]), 16))
		lastColon = false
	}
	return b.String()
}
