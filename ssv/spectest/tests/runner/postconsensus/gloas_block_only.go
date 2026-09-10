package postconsensus

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasBlockOnlyQuorum tests that a Gloas (ePBS) self-build proposer finalizes on a block-only
// post-consensus packet (SIP #94 §4/§6). The block root is the required entry and the §6 envelope root is
// optional, so a packet carrying only the block — what a run sees when the envelope partial-sigs never
// reach quorum — still submits the block and finishes the duty. The missing optional root is not an error,
// and with no envelope quorum there is no reveal to publish: the block alone is broadcast. This pins the
// liveness half of the fold — a stalled envelope round must not hold back the block.
func GloasBlockOnlyQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas

	return tests.NewMultiMsgProcessingSpecTest(
		"post consensus gloas block only",
		testdoc.PostConsensusGloasBlockOnlyDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name: "proposer block-only quorum finalizes",
				Runner: decideRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDutyV(version),
					testingutils.TestProposerConsensusDataV(version),
				),
				Duty: testingutils.TestingProposerDutyV(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerBlockOnlyMsgV(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerBlockOnlyMsgV(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerBlockOnlyMsgV(ks.Shares[3], 3, version))),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
				},
				DontStartDuty: true,
			},
		},
		ks,
	)
}
