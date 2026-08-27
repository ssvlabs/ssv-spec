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
// real Gloas block decision, recording the decided block root into a store the envelope runner shares,
// and the envelope flow then runs against that recorded root. Other envelope vectors seed the store
// with fixture literals; here the root the envelope duty consumes is the one §4 actually recorded, so
// any drift between the proposer's recording and the envelope value check's expectation fails the flow.
func RootLinkage() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	slot := phase0.Slot(testingutils.TestingDutySlotGloas)

	// Shared §4→§6 store, empty until the proposer decides.
	sharedRoots := ssv.ProposedBlockRoots{}

	// Run the §4 proposer to a decision for the slot; ProcessConsensus records the decided block root.
	proposerRunner := testingutils.ProposerRunner(ks).(*ssv.ProposerRunner)
	proposerRunner.ProposedBlockRoots = sharedRoots
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
	if root, ok := sharedRoots.Get(slot); !ok || root != testingutils.TestingProposedGloasBlockRoot(slot) {
		panic("proposer run did not record the decided block root")
	}

	// The envelope runner reads the recorded root; its value check still holds the fixture-seeded
	// store, so the two sides must agree for the flow to decide and publish.
	envelopeRunner := testingutils.EnvelopeProposerRunner(ks).(*ssv.EnvelopeProposerRunner)
	envelopeRunner.ProposedBlockRoots = sharedRoots

	return &tests.MsgProcessingSpecTest{
		Name:          "envelope proposer root linkage",
		Documentation: testdoc.EnvelopeProposerRootLinkageDoc,
		Runner:        envelopeRunner,
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
