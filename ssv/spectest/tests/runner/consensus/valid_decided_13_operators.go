package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidDecided13Operators tests a valid decided value (13 operators)
func ValidDecided13Operators() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing13SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid decided 13 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:                    "sync committee contribution",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeContributionConsensusDataByts, ks, types.BNRoleSyncCommitteeContribution),
				PostDutyRunnerStateRoot: "1e3dcd6dce0741033c1ad16663f4aa153873a224acb2967ad8a3605a5ea9b248",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  testingutils.SyncCommitteeRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeConsensusDataByts, ks, types.BNRoleSyncCommittee),
				PostDutyRunnerStateRoot: "462a0b4a0b38c02de8632af3d13cf9f2431e284e1bb044b26576acd73b4d4b9b",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    testingutils.TestingAggregatorDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestAggregatorConsensusDataByts, ks, types.BNRoleAggregator),
				PostDutyRunnerStateRoot: "a21736cd3af033159408afb5deb0356d3b5eae28de9414e15627a687ade99591",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestProposerConsensusDataByts, ks, types.BNRoleProposer),
				PostDutyRunnerStateRoot: "a86d337f76c522465dab459329ba00567caaca90006af683310f0fff555b3d92",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "attester",
				Runner:                  testingutils.AttesterRunner(ks),
				Duty:                    testingutils.TestingAttesterDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestAttesterConsensusDataByts, ks, types.BNRoleAttester),
				PostDutyRunnerStateRoot: "dff7b01202276bce1cf3c5b376e43361a2e8b28dd2f09bd5c3327635c58175d3",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
