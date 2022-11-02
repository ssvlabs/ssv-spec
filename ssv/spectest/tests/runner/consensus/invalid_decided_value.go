package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// InvalidDecidedValue tests an invalid decided value
func InvalidDecidedValue() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedError := "failed processing consensus message: decided ConsensusData invalid: decided value is invalid: duty invalid: wrong beacon role type"
	consensusData := func(role types.BeaconRole) *qbft.Data {
		cd := &types.ConsensusData{
			Duty: &types.Duty{
				Type:                    100,
				PubKey:                  testingutils.TestingValidatorPubKey,
				Slot:                    testingutils.TestingDutySlot,
				ValidatorIndex:          testingutils.TestingValidatorIndex,
				CommitteeIndex:          3,
				CommitteesAtSlot:        36,
				CommitteeLength:         128,
				ValidatorCommitteeIndex: 11,
			},
		}
		byts, _ := cd.MarshalSSZ()
		r, _ := cd.HashTreeRoot()
		return &qbft.Data{
			Root:   r,
			Source: byts,
		}
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "decided duty wrong role",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), types.PartialContributionProofSignatureMsgType),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2), types.PartialContributionProofSignatureMsgType),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3), types.PartialContributionProofSignatureMsgType),

					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							&qbft.Message{
								Height: qbft.FirstHeight,
								Round:  qbft.FirstRound,
								Input:  consensusData(types.BNRoleSyncCommitteeContribution),
							}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "1c72d49edf628f3d69bb0fcc2f4e2e0f96a593ebe20916d79bc885c2854fc6fc",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedError,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommittee(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							&qbft.Message{
								Height: qbft.FirstHeight,
								Round:  qbft.FirstRound,
								Input:  consensusData(types.BNRoleSyncCommittee),
							}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "971bbb10725791ebfd0a06f9817797feac2b81b17babb538b3c3f34421341cda",
				OutputMessages:          []*ssv.SignedPartialSignature{},
				ExpectedError:           "failed processing consensus message: failed to process consensus msg: could not process msg: invalid decided msg: decided value invalid: duty invalid: wrong beacon role type",
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), types.PartialSelectionProofSignatureMsgType),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2), types.PartialSelectionProofSignatureMsgType),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3), types.PartialSelectionProofSignatureMsgType),

					testingutils.SSVMsgAggregator(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							&qbft.Message{
								Height: qbft.FirstHeight,
								Round:  qbft.FirstRound,
								Input:  consensusData(types.BNRoleAggregator),
							}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "4ce7b608727db28f7af80765a8d98d1eced06f5bcc61cd5846b4ac15b1452cca",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedError,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1), types.PartialRandaoSignatureMsgType),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[2], ks.Shares[2], 2, 2), types.PartialRandaoSignatureMsgType),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[3], ks.Shares[3], 3, 3), types.PartialRandaoSignatureMsgType),

					testingutils.SSVMsgProposer(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							&qbft.Message{
								Height: qbft.FirstHeight,
								Round:  qbft.FirstRound,
								Input:  consensusData(types.BNRoleProposer),
							}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "374c4bd669d107076ef709c60c3e084fe3b27024a2cf1e3dbf59fc33c8dd6fb8",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
				ExpectedError: expectedError,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAttester(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							&qbft.Message{
								Height: qbft.FirstHeight,
								Round:  qbft.FirstRound,
								Input:  consensusData(types.BNRoleAttester),
							}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "d2ff2fad2e2b99af0f3e1f0378384ea9e89815d4dc4ca93f8754583a4a019b07",
				OutputMessages:          []*ssv.SignedPartialSignature{},
				ExpectedError:           "failed processing consensus message: failed to process consensus msg: could not process msg: invalid decided msg: decided value invalid: duty invalid: wrong beacon role type",
			},
		},
	}
}
