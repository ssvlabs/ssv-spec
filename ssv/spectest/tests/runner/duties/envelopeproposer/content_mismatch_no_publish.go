package envelopeproposer

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ContentMismatchNoPublish tests publish-by-content-match for a non-builder operator (SIP #94 §6): at the
// non-builder slot this operator's beacon node did not build the §4-decided block
// (GetBlindedExecutionPayloadEnvelope errors), so it disseminates nothing and only signs the envelope a
// builder operator disseminated. The signing round reaches quorum and the duty finishes, but this
// operator publishes nothing — only the builder operator holds the envelope body that blinds to the
// selected value.
func ContentMismatchNoPublish() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	slot := phase0.Slot(testingutils.TestingEnvelopeNonBuilderSlot)

	return &tests.MsgProcessingSpecTest{
		Name:          "envelope proposer non-builder no publish",
		Documentation: testdoc.EnvelopeProposerContentMismatchNoPublishDoc,
		Runner:        testingutils.EnvelopeProposerRunner(ks),
		Duty:          testingutils.TestingEnvelopeProposerNonBuilderDuty(),
		Messages: []*types.SignedSSVMessage{
			// The builder operator (op 2) disseminates the blinded envelope; this operator selects and signs it.
			testingutils.SignedEnvelopeDisseminationSSVMsg(ks, 2, slot),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PreConsensusEnvelopeMsgForSlot(ks.Shares[1], 1, slot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PreConsensusEnvelopeMsgForSlot(ks.Shares[2], 2, slot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PreConsensusEnvelopeMsgForSlot(ks.Shares[3], 3, slot))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusEnvelopeMsgForSlot(ks.Shares[1], 1, slot), // signs on receiving the binding dissemination
		},
		BeaconBroadcastedRoots: []string{}, // not the builder operator — no publish
	}
}
