package ssv

// BuilderEntry is one cluster-configured direct-builder connection for the builder-request-auth extension
// (SIP #94 §5). Data is the token agreed with the builder out of band, signed byte-for-byte; when empty it
// defaults to the UTF-8 bytes of URL. Entries MUST be byte-identical across all operators, and SSV caps
// them at MaxBuilderEntries per validator.
type BuilderEntry struct {
	Data []byte
	URL  string
}

// MaxBuilderEntries is SSV's per-validator cap on configured builder entries (SIP #94 §5) — a sub-cap of
// the beacon-API's MAX_BUILDER_ENTRIES (64) that bounds the §7 message-validation budget.
const MaxBuilderEntries = 8

// AuthData returns the bytes signed for this entry: the configured Data, or the URL bytes when Data is
// empty (SIP #94 §5). A zero-length result is invalid and is skipped by the auth round.
func (e BuilderEntry) AuthData() []byte {
	if len(e.Data) > 0 {
		return e.Data
	}
	return []byte(e.URL)
}
