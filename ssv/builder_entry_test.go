package ssv_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/ssv"
)

// TestBuilderEntryAuthData pins the builder-request-auth data resolution (SIP #94 §5): a configured Data is
// signed byte-for-byte, and an empty Data falls back to the builder URL's hostname per the upstream
// get_default_auth_data rule (builder-specs #168). The URL cases mirror the upstream vectors verbatim.
func TestBuilderEntryAuthData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		url  string
		want string
	}{
		// Configured Data wins and is signed verbatim; URL is ignored.
		{"configured data wins", []byte("builder-auth-token"), "https://ignored.example", "builder-auth-token"},

		// Upstream get_default_auth_data vectors (empty Data → hostname default).
		{"plain host", nil, "https://builder.example.com/", "builder.example.com"},
		{"uppercase scheme, port, path and query dropped", nil, "HTTPS://Builder.Example.com:443/bids?x=1", "builder.example.com"},
		{"explicit port dropped", nil, "https://builder.example.com:8080", "builder.example.com"},
		{"userinfo dropped", nil, "https://user:pw@builder.example.com/", "builder.example.com"},
		{"ipv4 literal", nil, "https://10.0.0.5:18550/eth/v1/builder", "10.0.0.5"},
		{"ipv6 literal compressed", nil, "https://[0:0:0:0:0:0:0:1]:8443/", "[::1]"},
		{"ipv4-mapped ipv6 stays hex-only", nil, "https://[::ffff:192.0.2.1]/", "[::ffff:c000:201]"},

		// Case-folding boundary between Go and upstream Python (SIP #94 §5, builder-specs #168): both fold
		// U+212A (Kelvin sign) to ASCII "k" and accept it, but only Go folds U+0130 (İ) to ASCII — Python folds
		// it to "i" + U+0307, fails isascii() and rejects, so U+0130 is rejected to keep the auth root aligned.
		{"U+212A (Kelvin) folds to ascii, accepted", nil, "https://\u212Aevin.example/", "kevin.example"},
		{"U+0130 (dotted I) rejected to match upstream", nil, "https://\u0130stanbul.example/", ""},

		// No derivable ASCII hostname → no default (the auth round skips the entry).
		{"no scheme yields no host", nil, "builder.example.com", ""},
		{"non-ascii host must be punycode", nil, "https://exämple.com/", ""},
		{"empty entry", nil, "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ssv.BuilderEntry{Data: tc.data, URL: tc.url}.AuthData()
			require.Equal(t, tc.want, string(got))
		})
	}
}
