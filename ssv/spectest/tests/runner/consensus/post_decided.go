package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests a valid commit msg after returned decided already
func PostDecided() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeContributionConsensusDataByts, ks, types.BNRoleSyncCommitteeContribution),
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeContributionMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
						}), nil)),
				PostDutyRunnerStateRoot: "62aaf957e5fb61294d6ed18c24b9213badc850a0bf5f3a35787c341f83c29aed",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeConsensusDataByts, ks, types.BNRoleSyncCommittee),
					testingutils.SSVMsgSyncCommittee(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
						}), nil)),
				PostDutyRunnerStateRoot: "65253eaf02a6874ebbaa1c91e311054c68e42f5e26c3e43774d28a6b4078c5cb",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestAggregatorConsensusDataByts, ks, types.BNRoleAggregator),
					testingutils.SSVMsgAggregator(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AggregatorMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
						}), nil)),
				PostDutyRunnerStateRoot: "dd2c5fb1a4dc2bab521cb6155129dafb553c1e92f1285a76500d33f876571855",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestProposerConsensusDataByts, ks, types.BNRoleProposer),
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
						}), nil)),
				PostDutyRunnerStateRoot: "b1ae6786aa5e442e65a64ed428c9c13a35d277fcc5cceaa700c1305d0e53839e",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
					testingutils.PostConsensusProposerMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestAttesterConsensusDataByts, ks, types.BNRoleAttester),
					testingutils.SSVMsgAttester(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AttesterMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
						}), nil)),
				PostDutyRunnerStateRoot: "eb592a006bdb01d58d159af49a0be32d3f2092de5274426391906726e9728e8c",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(testingutils.Testing4SharesSet().Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
