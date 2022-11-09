package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// FutureDecided tests a running instance at FirstHeight, then processing a decided msg from height 10 and returning decided but doesn't move to post consensus as it's not the same instance decided
func FutureDecided() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	inputData := &qbft.Data{
		Root:   testingutils.TestAttesterConsensusDataRoot,
		Source: testingutils.TestAttesterConsensusDataByts,
	}
	errStr := "failed processing consensus message: decided wrong instance"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus future decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), types.PartialContributionProofSignatureMsgType),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2), types.PartialContributionProofSignatureMsgType),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3), types.PartialContributionProofSignatureMsgType),
					testingutils.SSVMsgSyncCommitteeContribution(testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							Height: 2,
							Round:  qbft.FirstRound,
						}, inputData), nil, types.DecidedMsgType),
				},
				PostDutyRunnerStateRoot: "5bc67309b58de15eaeb3febf7ad9474362ad20c63fc3af697c2497cc23187fdf",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommittee(testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							Height: 2,
							Round:  qbft.FirstRound,
						}, inputData), nil, types.DecidedMsgType),
				},
				PostDutyRunnerStateRoot: "4b07a9a4c0e2d7b618dd350733d32bed8ef946fc2d8144e4f873c245ac47d99a",
				OutputMessages:          []*ssv.SignedPartialSignatures{},
				ExpectedError:           errStr,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), types.PartialSelectionProofSignatureMsgType),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2), types.PartialSelectionProofSignatureMsgType),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3), types.PartialSelectionProofSignatureMsgType),
					testingutils.SSVMsgAggregator(testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							Height: 2,
							Round:  qbft.FirstRound,
						}, inputData), nil, types.DecidedMsgType),
				},
				PostDutyRunnerStateRoot: "9c33e371726895b344d612bf342f96004041362fbf4c942f51606579abcf9aef",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), types.PartialRandaoSignatureMsgType),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[2], 2), types.PartialRandaoSignatureMsgType),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[3], 3), types.PartialRandaoSignatureMsgType),
					testingutils.SSVMsgProposer(testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							Height: 2,
							Round:  qbft.FirstRound,
						}, inputData), nil, types.DecidedMsgType),
				},
				PostDutyRunnerStateRoot: "fe3c2744650f62512df4c624b38daf7326a804ea2d72ea794905fec91bb5bd32",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAttester(testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							Height: 2,
							Round:  qbft.FirstRound,
						}, inputData), nil, types.DecidedMsgType),
				},
				PostDutyRunnerStateRoot: "20d49680b13605f2f0769b707e7d3170180552840400ef5138cc58ed87ee9538",
				OutputMessages:          []*ssv.SignedPartialSignatures{},
				ExpectedError:           errStr,
			},
		},
	}
}
