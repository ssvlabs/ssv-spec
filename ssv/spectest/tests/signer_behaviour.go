package tests

import (
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignerBehaviour configures non-default testing signer responses. Like BeaconNodeBehaviour it is part of the
// vector JSON, so an implementation running the vector configures its own signer mock the same way.
type SignerBehaviour struct {
	// FailEnvelopeSigning makes signing under DomainBeaconBuilder (the §6 envelope root) fail (SIP #94 §4).
	FailEnvelopeSigning bool `json:",omitempty"`
}

// Apply configures km; a nil behaviour leaves the defaults.
func (b *SignerBehaviour) Apply(km *testingutils.TestingKeyManager) {
	if b == nil {
		return
	}
	km.SetFailEnvelopeSigning(b.FailEnvelopeSigning)
}
