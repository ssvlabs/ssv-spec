package runner

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasProposerNotBuilderOperator tests that only the builder operator publishes the §6 reveal (SIP #94 §6).
// The cluster decides a self-build block this operator did not produce, so it signs both post-consensus roots
// and submits the block, but publishes no envelope: its own produce was either a self-build of another block,
// whose envelope is not the decided block's, or an external bid, which has no envelope.
func GloasProposerNotBuilderOperator() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	// The same run for both sub-tests; they differ in what this operator's own beacon node produced.
	test := func(name string, beaconNode *tests.BeaconNodeBehaviour) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   name,
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer), // pre-consensus + consensus
				[]*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
				}...,
			),
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			// The block only: the envelope root reaches quorum too, but this operator has no reveal to publish.
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
			},
			BeaconNode: beaconNode,
		}
	}

	return tests.NewMultiMsgProcessingSpecTest(
		"gloas proposer not builder operator",
		testdoc.GloasProposerNotBuilderOperatorDoc,
		[]*tests.MsgProcessingSpecTest{
			test("produced another self-build block", &tests.BeaconNodeBehaviour{ProducesOtherSelfBuildBlock: true}),
			test("produced an external bid", &tests.BeaconNodeBehaviour{ProducesExternalBid: true}),
		},
		ks,
	)
}
