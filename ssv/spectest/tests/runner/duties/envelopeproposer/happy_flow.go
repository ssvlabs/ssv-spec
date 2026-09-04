package envelopeproposer

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// HappyFlow tests the full §6 envelope flow for the builder operator (SIP #94 §6): the duty starts for a
// slot whose §4-decided block is recorded, this operator's beacon node built that block, so it
// disseminates its blinded envelope and threshold-signs it in one round; the signing round reaches
// quorum and — since it holds the produced envelope that blinds to the selected one — it publishes the
// reveal. There is no consensus phase.
func HappyFlow() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	slot := testingutils.TestingEnvelopeProposerDuty().Slot

	return &tests.MsgProcessingSpecTest{
		Name:          "envelope proposer happy flow",
		Documentation: testdoc.EnvelopeProposerHappyFlowDoc,
		Runner:        testingutils.EnvelopeProposerRunner(ks),
		Duty:          testingutils.TestingEnvelopeProposerDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PreConsensusEnvelopeMsg(ks.Shares[1], 1))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PreConsensusEnvelopeMsg(ks.Shares[2], 2))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PreConsensusEnvelopeMsg(ks.Shares[3], 3))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusEnvelopeMsg(ks.Shares[1], 1), // disseminates and signs on duty start
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingBlindedExecutionPayloadEnvelope(slot)),
		},
	}
}
