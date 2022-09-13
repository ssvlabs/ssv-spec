package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoRunningDuty tests a valid proposal msg before duty starts
func NoRunningDuty() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus no running duty",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeContributionMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "74deeed09b4bda98d0ed5f91d0a46378a664e05c80a2b1aa3e2a9feec56dcf73",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestSyncCommitteeConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "a024bcb3737480ad2801e5a0db8890e4cc90e17cc59721f4ccc7cf8816a4d2f6",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AggregatorMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestAggregatorConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "e586a176b2ab97ce6ba3b28cc205085f285ebe926533eae9f2896e3edb029afe",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestProposerConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "e24a4beacc971c1a27b3f4fe2bc04f8b6a3b290a17bdf25e60e88ee3d2085bd3",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AttesterMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestAttesterConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "ed2ab5ff75817cc62dfb7c0bc1e093f51633d2c58d2c0aa0bf76eb56d20add9e",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
		},
	}
}
