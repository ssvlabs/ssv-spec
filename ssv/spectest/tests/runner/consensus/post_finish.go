package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostFinish tests a valid commit msg after runner finished
func PostFinish() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	finishRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances[0] = r.GetBaseRunner().State.RunningInstance
		r.GetBaseRunner().QBFTController.Height = qbft.FirstHeight
		r.GetBaseRunner().State.Finished = true
		return r
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post finish",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: finishRunner(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionDuty),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestSyncCommitteeContributionConsensusDataRoot,
							Source: testingutils.TestSyncCommitteeContributionConsensusDataByts,
						}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "a7749c274f3b3add80d9e6444c3c5cd720f458e9ce41ef36789dd4ca1c79eac5",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				DontStartDuty: true,
				ExpectedError: "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "sync committee",
				Runner: finishRunner(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommittee(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestSyncCommitteeConsensusDataRoot,
							Source: testingutils.TestSyncCommitteeConsensusDataByts,
						}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "90a493699ebbd1558cd6ba8ea5e8045a7ab980c7e4720889310956792539ee50",
				OutputMessages:          []*ssv.SignedPartialSignatures{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "aggregator",
				Runner: finishRunner(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDuty),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAggregator(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestAggregatorConsensusDataRoot,
							Source: testingutils.TestAggregatorConsensusDataByts,
						}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "69400f37b9724938c62791cfff2411cfc5722df878630a803827491b6ec74b25",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				DontStartDuty: true,
				ExpectedError: "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "proposer",
				Runner: finishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDuty),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestProposerConsensusDataRoot,
							Source: testingutils.TestProposerConsensusDataByts,
						}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "265d356d8f3fc6c53ca975b99d6137047ba9b04a672b8ca67f5878edc84c466e",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
				DontStartDuty: true,
				ExpectedError: "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "attester",
				Runner: finishRunner(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAttester(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestAttesterConsensusDataRoot,
							Source: testingutils.TestAttesterConsensusDataByts,
						}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "2cab1cbf84b00697333c01973c5ee7cce950ec34762d0811af6cad3f6b87512b",
				OutputMessages:          []*ssv.SignedPartialSignatures{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
		},
	}
}
