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
		byts, _ := cd.MarshalSSZ()
		return byts
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
								Input: &qbft.Data{
									Root:   [32]byte{},
									Source: consensusDataByts(types.BNRoleSyncCommitteeContribution),
								},
							}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "c657220938d48e0573e9874b84b29cb92c7d581fca27c28d2a34cc3913c3ffbb",
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
								Input: &qbft.Data{
									Root:   [32]byte{},
									Source: consensusDataByts(types.BNRoleSyncCommittee),
								},
							}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "ca7c1b5bb6a1b2b5d486da30bdef8a96db109cbe5691d1191a0671eaaafb5cf0",
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
								Input: &qbft.Data{
									Root:   [32]byte{},
									Source: consensusDataByts(types.BNRoleAggregator),
								},
							}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "e5df3344d76b2e83fabbff92998d5999daff2747f276b1fbc222b6c960750b93",
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
								Input: &qbft.Data{
									Root:   [32]byte{},
									Source: consensusDataByts(types.BNRoleProposer),
								},
							}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "f9023766285ee75a894437f41dbe7d40b6d23922df60353a6c0125d56c4f9f09",
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
								Input: &qbft.Data{
									Root:   [32]byte{},
									Source: consensusDataByts(types.BNRoleAttester),
								},
							}), nil, types.ConsensusCommitMsgType),
				},
				PostDutyRunnerStateRoot: "0a87eee70b6ee2583dd414d6f07f6f5c433975409896dafb51628f5e393a7458",
				OutputMessages:          []*ssv.SignedPartialSignature{},
				ExpectedError:           "failed processing consensus message: failed to process consensus msg: could not process msg: invalid decided msg: decided value invalid: duty invalid: wrong beacon role type",
			},
		},
	}
}
