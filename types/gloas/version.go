package gloas

import (
	"github.com/attestantio/go-eth2-client/spec"
)

// DataVersionGloas is a placeholder for the Gloas beacon data version: until go-eth2-client
// defines it (its latest is Fulu), Gloas is slotted immediately after Fulu. The upstream enum is
// iota-ordered, so this value matches what upstream will assign. Remove and reconcile once
// upstream ships a real spec.DataVersionGloas. Note DataVersionGloas.String() returns "unknown" —
// the upstream string/JSON tables aren't extended.
const DataVersionGloas = spec.DataVersionFulu + 1
