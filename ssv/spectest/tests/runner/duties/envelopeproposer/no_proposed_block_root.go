package envelopeproposer

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoProposedBlockRoot tests that the envelope duty takes no action for a slot with no recorded §4-decided
// block (SIP #94 §6): the operator has not decided §4, so it cannot be the builder operator and cannot
// bind a dissemination. It disseminates nothing, signs nothing, and stays running until a decision
// arrives — this is a silent wait, not an error.
func NoProposedBlockRoot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	duty := testingutils.TestingEnvelopeProposerDuty()
	duty.Slot = testingutils.TestingDutySlotGloasInvalid // no recorded block for this slot

	return &tests.MsgProcessingSpecTest{
		Name:                   "envelope proposer no proposed block",
		Documentation:          testdoc.EnvelopeProposerNoProposedBlockRootDoc,
		Runner:                 testingutils.EnvelopeProposerRunner(ks),
		Duty:                   duty,
		Messages:               []*types.SignedSSVMessage{},
		OutputMessages:         []*types.PartialSignatureMessages{},
		BeaconBroadcastedRoots: []string{},
	}
}
