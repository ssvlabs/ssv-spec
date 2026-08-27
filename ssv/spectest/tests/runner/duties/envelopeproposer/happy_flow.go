package envelopeproposer

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// HappyFlow tests the full §6 envelope flow (SIP #94 §6): the duty starts for a slot whose §4-decided
// block root is recorded, consensus decides the operator's own produced blinded envelope, post-consensus
// reaches quorum, and — since this operator produced the decided envelope — it publishes the signed
// envelope.
func HappyFlow() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	slot := testingutils.TestingEnvelopeProposerDuty().Slot

	return &tests.MsgProcessingSpecTest{
		Name:          "envelope proposer happy flow",
		Documentation: testdoc.EnvelopeProposerHappyFlowDoc,
		Runner:        testingutils.EnvelopeProposerRunner(ks),
		Duty:          testingutils.TestingEnvelopeProposerDuty(),
		Messages: append(
			testingutils.SSVDecidingMsgsForEnvelopeProposer(slot, ks),
			[]*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PostConsensusEnvelopeProposerMsg(ks.Shares[1], 1))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PostConsensusEnvelopeProposerMsg(ks.Shares[2], 2))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PostConsensusEnvelopeProposerMsg(ks.Shares[3], 3))),
			}...,
		),
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PostConsensusEnvelopeProposerMsg(ks.Shares[1], 1), // broadcasts when consensus decides
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBlindedExecutionPayloadEnvelope(ks, slot)),
		},
	}
}
