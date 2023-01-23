package consensus

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
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
				PostDutyRunnerStateRoot: "b0b770dfb03fb2b2530f0e3d1ea269639c54ee03f4a75d18a18e8c8324b8f637",
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
				PostDutyRunnerStateRoot: "7a201ae85b215a1af4732a2d519a3d28d2848aac1724cd4d24a103b05d02b683",
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
				PostDutyRunnerStateRoot: "0df459cdbb965e07dc63493e656f8e0086b65bd4e7b4d922857e63097bd3f33a",
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
				PostDutyRunnerStateRoot: "e9e011cec6219def6bcd545fbb8201b12ae38e5f4eb6713a90745145e0b2710e",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
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
				PostDutyRunnerStateRoot: "cad75704dcd80b4687362fbd3a3c870e6c7abd65d874a6986163ab8e2f08b6b5",
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
				PostDutyRunnerStateRoot: "ae50555585c7f0a4d8b3d1a6c9138bb00d8b0ebdb62d6b792f7187eeb5a6e4b0",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           errStr,
			},
		},
	}
}
