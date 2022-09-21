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
				PostDutyRunnerStateRoot: "95d33f7387e349b91f7c257028d71f634d1558b468945030f1ee78d0682f5ad5",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedError,
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
				PostDutyRunnerStateRoot: "6f22c1792cc1a90ba96a3053588f2c6e13e13a510fc6e9c1ca99d6b140f521d1",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
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
				PostDutyRunnerStateRoot: "311e08d705481d31207d79d81200529c6c553fb11ab8f7eff436e4c3cd0f871b",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedError,
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
				PostDutyRunnerStateRoot: "1bd4b88c1748ff0a54d2e86b0376fea7483e15b586211be13ca5f10ee5343ea9",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
				ExpectedError: expectedError,
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
				PostDutyRunnerStateRoot: "c07a579655fba0e74e283289e91ade60157a54c23b31f1a18774ac29baa08af9",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
		},
	}
}
