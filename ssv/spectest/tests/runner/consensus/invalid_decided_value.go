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
				PostDutyRunnerStateRoot: "c657220938d48e0573e9874b84b29cb92c7d581fca27c28d2a34cc3913c3ffbb",
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
				PostDutyRunnerStateRoot: "4ecc5fa28895cda24f3c9489c891c69431ad1c9aedf66d27f87528494b313e19",
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
				PostDutyRunnerStateRoot: "e5df3344d76b2e83fabbff92998d5999daff2747f276b1fbc222b6c960750b93",
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
				PostDutyRunnerStateRoot: "f9023766285ee75a894437f41dbe7d40b6d23922df60353a6c0125d56c4f9f09",
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
				PostDutyRunnerStateRoot: "d00a1795c598a53727f7b683d54b2299289dd879c82198123556537e781fcede",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "failed processing consensus message: failed to process consensus msg: could not process msg: invalid decided msg: decided value invalid: duty invalid: wrong beacon role type",
			},
		},
	}
}
