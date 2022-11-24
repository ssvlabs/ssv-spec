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

	startInstance := func(r ssv.Runner, value []byte) ssv.Runner {
		if err := r.GetBaseRunner().QBFTController.StartNewInstance(value); err != nil {
			panic(err.Error())
		}
		return r
	}

	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus no running duty",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee contribution",
				Runner: startInstance(
					testingutils.SyncCommitteeContributionRunner(ks),
					testingutils.TestSyncCommitteeContributionConsensusDataByts,
				),
				Duty: testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeContributionMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "2357948c78fb602da32ba81f12aba09fa798ac5ee19701c03ad0b2b73af6254f",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name: "sync committee",
				Runner: startInstance(
					testingutils.SyncCommitteeRunner(ks),
					testingutils.TestSyncCommitteeConsensusDataByts,
				),
				Duty: testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestSyncCommitteeConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "169afe8e7115d7003acb49177558abb440ce57651eecaf37b05dd5693c4cecda",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name: "aggregator",
				Runner: startInstance(
					testingutils.AggregatorRunner(ks),
					testingutils.TestAggregatorConsensusDataByts,
				),
				Duty: testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AggregatorMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestAggregatorConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "56b533bf403e9cb4d19abc39e7805bd0fa8579321baacd583d1a187d86df36d2",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name: "proposer",
				Runner: startInstance(
					testingutils.ProposerRunner(ks),
					testingutils.TestProposerConsensusDataByts,
				),
				Duty: testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestProposerConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "840a04743b2bbadfef80a512e6727c7adeaa1a193a0e743810d0061cde35ad3e",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name: "attester",
				Runner: startInstance(
					testingutils.AttesterRunner(ks),
					testingutils.TestAttesterConsensusDataByts,
				),
				Duty: testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AttesterMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestAttesterConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "c9eb871dc586f9e23f645667eb8f5b21ecc7308af7557c1a25d880f9bb6dedaf",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
		},
	}
}
