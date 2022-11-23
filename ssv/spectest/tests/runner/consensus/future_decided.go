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

	getID := func(role types.BeaconRole) []byte {
		ret := types.NewMsgID(testingutils.TestingValidatorPubKey[:], role)
		return ret[:]
	}

	errStr := "failed processing consensus message: decided wrong instance"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus future decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					testingutils.SSVMsgSyncCommitteeContribution(testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     2,
							Round:      qbft.FirstRound,
							Identifier: getID(types.BNRoleSyncCommitteeContribution),
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}), nil),
				},
				PostDutyRunnerStateRoot: "697f350d97af8b5f91c3059f4a77780e0297c7be0c686ab5d14eb6c6e71d9c28",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     2,
							Round:      qbft.FirstRound,
							Identifier: getID(types.BNRoleSyncCommittee),
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}), nil),
				},
				PostDutyRunnerStateRoot: "3e6afc06b9db2e8b7b678655783fc12c17b5fe078d5f90f50b453c910238fd14",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           errStr,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					testingutils.SSVMsgAggregator(testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     2,
							Round:      qbft.FirstRound,
							Identifier: getID(types.BNRoleAggregator),
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}), nil),
				},
				PostDutyRunnerStateRoot: "20384a054b86cbb608da85f5b9c88a4d09ed592702d7a91326f0412d0bc403a3",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgProposer(testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     2,
							Round:      qbft.FirstRound,
							Identifier: getID(types.BNRoleProposer),
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}), nil),
				},
				PostDutyRunnerStateRoot: "c6d1cba30aff597d2bcc3f8236dfd2ad09953f3f62f8edc6932a1889fea2f11e",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     2,
							Round:      qbft.FirstRound,
							Identifier: getID(types.BNRoleAttester),
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}), nil),
				},
				PostDutyRunnerStateRoot: "99284ed3cf0baf82d6e0ce6ce062fa9071a6607001e2d098e0a5920449f51d1d",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           errStr,
			},
		},
	}
}
