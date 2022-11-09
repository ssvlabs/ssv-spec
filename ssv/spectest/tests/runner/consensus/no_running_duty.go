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
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestSyncCommitteeContributionConsensusDataRoot,
							Source: testingutils.TestSyncCommitteeContributionConsensusDataByts,
						}), nil, types.ConsensusProposeMsgType),
				},
				PostDutyRunnerStateRoot: "89fd84d043d6a281f5e332059fc08d77c5dcdcbc5fd30a2c4acc261d0463e32e",
				OutputMessages:          []*ssv.SignedPartialSignatures{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommittee(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestSyncCommitteeConsensusDataRoot,
							Source: testingutils.TestSyncCommitteeConsensusDataByts,
						}), nil, types.ConsensusProposeMsgType),
				},
				PostDutyRunnerStateRoot: "8ab1f7abacc7897fe9558050764125348186ee5f21879fb6b004f1c51cac19db",
				OutputMessages:          []*ssv.SignedPartialSignatures{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAggregator(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestAggregatorConsensusDataRoot,
							Source: testingutils.TestAggregatorConsensusDataByts,
						}), nil, types.ConsensusProposeMsgType),
				},
				PostDutyRunnerStateRoot: "3aa9d352b7e5cd28be240ef9fa6614ccb97c74917790ca11fc3aa0d5a84f5ce9",
				OutputMessages:          []*ssv.SignedPartialSignatures{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestProposerConsensusDataRoot,
							Source: testingutils.TestProposerConsensusDataByts,
						}), nil, types.ConsensusProposeMsgType),
				},
				PostDutyRunnerStateRoot: "017ba7ecaf768de24235fd7b1b45fa6dade9f330190545dfc413672bdf5aea46",
				OutputMessages:          []*ssv.SignedPartialSignatures{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAttester(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestAttesterConsensusDataRoot,
							Source: testingutils.TestAttesterConsensusDataByts,
						}), nil, types.ConsensusProposeMsgType),
				},
				PostDutyRunnerStateRoot: "e747803aebf9b8c680b20ccb08adc932e44957df8d049b1c44dd0e4cdeb97818",
				OutputMessages:          []*ssv.SignedPartialSignatures{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
		},
	}
}
