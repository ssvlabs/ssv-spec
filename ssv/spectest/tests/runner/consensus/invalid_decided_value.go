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

	consensusDataByts := func(role types.BeaconRole) []byte {
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
		byts, _ := cd.Encode()
		return byts
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "decided duty wrong role",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),

					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							&qbft.Message{
								MsgType:    qbft.CommitMsgType,
								Height:     qbft.FirstHeight,
								Round:      qbft.FirstRound,
								Identifier: testingutils.SyncCommitteeContributionMsgID,
								Data:       testingutils.CommitDataBytes(consensusDataByts(types.BNRoleSyncCommitteeContribution)),
							}), nil),
				},
				PostDutyRunnerStateRoot: "f41289fab2fcfd8f4ddd5b87e7a49cef792003e74d2132558310677b9aebc6ae",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: "failed processing consensus message: failed to process consensus msg: could not process msg: invalid decided msg: decided value invalid: duty invalid: wrong beacon role type",
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							&qbft.Message{
								MsgType:    qbft.CommitMsgType,
								Height:     qbft.FirstHeight,
								Round:      qbft.FirstRound,
								Identifier: testingutils.SyncCommitteeMsgID,
								Data:       testingutils.CommitDataBytes(consensusDataByts(types.BNRoleSyncCommittee)),
							}), nil),
				},
				PostDutyRunnerStateRoot: "5ff863a832e766731418768336ccb6580bb3dd33c2a0436d0196ea6525a0749a",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "failed processing consensus message: failed to process consensus msg: could not process msg: invalid decided msg: decided value invalid: duty invalid: wrong beacon role type",
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),

					testingutils.SSVMsgAggregator(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							&qbft.Message{
								MsgType:    qbft.CommitMsgType,
								Height:     qbft.FirstHeight,
								Round:      qbft.FirstRound,
								Identifier: testingutils.AggregatorMsgID,
								Data:       testingutils.CommitDataBytes(consensusDataByts(types.BNRoleAggregator)),
							}), nil),
				},
				PostDutyRunnerStateRoot: "c32aa3355d684d473f15474129903aa2968068f57e3b9ccc1e8352372e796c1d",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: "failed processing consensus message: failed to process consensus msg: could not process msg: invalid decided msg: decided value invalid: duty invalid: wrong beacon role type",
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[3], ks.Shares[3], 3, 3)),

					testingutils.SSVMsgProposer(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							&qbft.Message{
								MsgType:    qbft.CommitMsgType,
								Height:     qbft.FirstHeight,
								Round:      qbft.FirstRound,
								Identifier: testingutils.ProposerMsgID,
								Data:       testingutils.CommitDataBytes(consensusDataByts(types.BNRoleProposer)),
							}), nil),
				},
				PostDutyRunnerStateRoot: "9b7ec2bce48fe2bb0582c04531d1fc797640515cd0d464c0a810a7390d9623dd",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
				ExpectedError: "failed processing consensus message: failed to process consensus msg: could not process msg: invalid decided msg: decided value invalid: duty invalid: wrong beacon role type",
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						testingutils.MultiSignQBFTMsg(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							&qbft.Message{
								MsgType:    qbft.CommitMsgType,
								Height:     qbft.FirstHeight,
								Round:      qbft.FirstRound,
								Identifier: testingutils.AttesterMsgID,
								Data:       testingutils.CommitDataBytes(consensusDataByts(types.BNRoleAttester)),
							}), nil),
				},
				PostDutyRunnerStateRoot: "e80945f69c9142e391550f605d18a1e39f423d5f865ef865421f07ec25c8b803",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "failed processing consensus message: failed to process consensus msg: could not process msg: invalid decided msg: decided value invalid: duty invalid: wrong beacon role type",
			},
		},
	}
}
