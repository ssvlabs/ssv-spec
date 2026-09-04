package envelopeproposer

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RootLinkage tests the §4→§6 linkage end to end (SIP #94 §6): a proposer runner is driven through a
// real Gloas block decision, recording the decided block facts into a store the envelope runner shares,
// and the builder-operator envelope flow then runs and binds against those recorded facts. Other
// envelope vectors seed the store with fixture facts; here the facts the envelope duty binds against are
// the ones §4 actually recorded, so any drift between the proposer's recording and the envelope runner's
// binding fails the flow.
func RootLinkage() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	slot := phase0.Slot(testingutils.TestingDutySlotGloas)

	// Shared §4→§6 store, empty until the proposer decides.
	sharedBlocks := ssv.ProposedBlocks{}

	// Run the §4 proposer to a decision for the slot; ProcessConsensus records the decided block facts.
	proposerRunner := testingutils.ProposerRunner(ks).(*ssv.ProposerRunner)
	proposerRunner.ProposedBlocks = sharedBlocks
	v := testingutils.BaseValidator(ks)
	v.DutyRunners[types.RoleProposer] = proposerRunner
	v.Network = proposerRunner.GetNetwork()
	if err := v.StartDuty(testingutils.TestingProposerDutyV(gloas.DataVersionGloas)); err != nil {
		panic(err.Error())
	}
	for _, msg := range testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(gloas.DataVersionGloas), ks, types.RoleProposer) {
		if err := v.ProcessMessage(msg); err != nil {
			panic(err.Error())
		}
	}
	if block, ok := sharedBlocks.Get(slot); !ok || block.BlockRoot != testingutils.TestingProposedGloasBlockRoot(slot) {
		panic("proposer run did not record the decided block")
	}

	// The envelope runner (builder operator) binds against the facts §4 recorded.
	envelopeRunner := testingutils.EnvelopeProposerRunner(ks).(*ssv.EnvelopeProposerRunner)
	envelopeRunner.ProposedBlocks = sharedBlocks

	return &tests.MsgProcessingSpecTest{
		Name:          "envelope proposer root linkage",
		Documentation: testdoc.EnvelopeProposerRootLinkageDoc,
		Runner:        envelopeRunner,
		Duty:          testingutils.TestingEnvelopeProposerDuty(),
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PreConsensusEnvelopeMsg(ks.Shares[1], 1))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PreConsensusEnvelopeMsg(ks.Shares[2], 2))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, testingutils.PreConsensusEnvelopeMsg(ks.Shares[3], 3))),
		},
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusEnvelopeMsg(ks.Shares[1], 1),
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingBlindedExecutionPayloadEnvelope(slot)),
		},
	}
}
