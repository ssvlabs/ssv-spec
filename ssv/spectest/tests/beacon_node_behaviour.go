package tests

import (
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BeaconNodeBehaviour configures non-default testing beacon node responses. It is part of the vector JSON, so an
// implementation running the vector configures its own beacon node mock the same way.
type BeaconNodeBehaviour struct {
	// ProducesExternalBid makes block production return a bare external-bid block and no envelope (SIP #94 §4).
	ProducesExternalBid bool `json:",omitempty"`
	// ProducesOtherSelfBuildBlock makes block production return a self-build block other than the one the vector's
	// messages decide, with its own envelope, so this operator is not the builder operator (SIP #94 §6).
	ProducesOtherSelfBuildBlock bool `json:",omitempty"`
	// FailBlockSubmit makes the block submit fail without recording the block (SIP #94 §6).
	FailBlockSubmit bool `json:",omitempty"`
	// WrongPayloadAttestationSlot makes the payload attestation data carry a slot other than the requested one
	// (SIP #94 §3).
	WrongPayloadAttestationSlot bool `json:",omitempty"`
	// FailProposerDutiesDependentRoot makes the proposer-duties dependent-root request fail (SIP #94 §5).
	FailProposerDutiesDependentRoot bool `json:",omitempty"`
}

// Apply configures bn; a nil behaviour leaves the defaults.
func (b *BeaconNodeBehaviour) Apply(bn *testingutils.TestingBeaconNode) {
	if b == nil {
		return
	}
	bn.SetProducesExternalBid(b.ProducesExternalBid)
	bn.SetProducesOtherSelfBuildBlock(b.ProducesOtherSelfBuildBlock)
	bn.SetFailBlockSubmit(b.FailBlockSubmit)
	bn.SetWrongPayloadAttestationSlot(b.WrongPayloadAttestationSlot)
	bn.SetFailProposerDutiesDependentRoot(b.FailProposerDutiesDependentRoot)
}
