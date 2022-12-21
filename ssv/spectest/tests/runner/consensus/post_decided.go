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
				PostDutyRunnerStateRoot: "b0e9c1e7093ff65f2ca95d576636e7cf47af119ea4a558ac868de331cd2d8140",
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
				PostDutyRunnerStateRoot: "03d9f0582c39a878689163387681d32a2b857ab7104ef343cf2e6d6a9fdbaaf5",
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
				PostDutyRunnerStateRoot: "d76feb4b969863f38af71df5c412348f4c12292e137fa307009155adcdcb9801",
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
				PostDutyRunnerStateRoot: "edd00237cbfeb826243c7aa6fe50498a2d9156063e399ac49f70adb867158a73",
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
				PostDutyRunnerStateRoot: "88b5f91f8a76ac4ca521854420a766f7209549d78049f4a89efa23c0fca989a4",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(testingutils.Testing4SharesSet().Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
