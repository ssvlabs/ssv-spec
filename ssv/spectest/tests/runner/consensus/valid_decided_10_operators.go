package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidDecided10Operators tests a valid decided value (10 operators)
func ValidDecided10Operators() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing10SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid decided 10 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:                    "sync committee contribution",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeContributionConsensusDataByts, ks, types.BNRoleSyncCommitteeContribution),
				PostDutyRunnerStateRoot: "5eccf1fa79184350d02d3533e1ca1c3ddaf30b7d15d34d01c96f0d730ea7afd9",
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
				PostDutyRunnerStateRoot: "356fb790ec403278fa705d17db18ee600f0615c824d6bafaa13c9a092bdd5c5c",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    testingutils.TestingAggregatorDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestAggregatorConsensusDataByts, ks, types.BNRoleAggregator),
				PostDutyRunnerStateRoot: "e49905b6eed09c752618ceb01898a0fc1974dd030943c850c1afabc061896f1d",
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
				PostDutyRunnerStateRoot: "426c00c2fb2b400f14305b444889f4546172bee6882b1c5c3f86b6534d1f2a0f",
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
				PostDutyRunnerStateRoot: "590123c6d684a75271e056cebf6e19ecb206332f5afee67f71519a347d4d4034",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
