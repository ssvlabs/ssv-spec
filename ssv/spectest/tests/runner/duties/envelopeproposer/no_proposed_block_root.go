package envelopeproposer

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoProposedBlockRoot tests that the envelope duty is rejected for a slot with no recorded §4-decided
// block root (SIP #94 §6): the duty only exists after the proposer decided, so nothing can be built
// or agreed without the root.
func NoProposedBlockRoot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	duty := testingutils.TestingEnvelopeProposerDuty()
	duty.Slot = testingutils.TestingDutySlotGloasInvalid // no recorded root for this slot

	return &tests.MsgProcessingSpecTest{
		Name:              "envelope proposer no proposed block root",
		Documentation:     testdoc.EnvelopeProposerNoProposedBlockRootDoc,
		Runner:            testingutils.EnvelopeProposerRunner(ks),
		Duty:              duty,
		Messages:          []*types.SignedSSVMessage{},
		OutputMessages:    []*types.PartialSignatureMessages{},
		ExpectedErrorCode: types.EnvelopeNoProposedBlockRootErrorCode,
	}
}
