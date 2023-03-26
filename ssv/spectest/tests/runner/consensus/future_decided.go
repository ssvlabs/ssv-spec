package consensus

import (
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FutureDecided tests a running instance at FirstHeight, then processing a decided msg from height 10 and returning decided but doesn't move to post consensus as it's not the same instance decided
func FutureDecided() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	getID := func(role types.BeaconRole) []byte {
		ret := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], role)
		return ret[:]
	}

	errStr := "failed processing consensus message: decided wrong instance"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus future decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							getID(types.BNRoleSyncCommitteeContribution),
						),
						nil,
					),
				},
				PostDutyRunnerStateRoot: "6e003cd912d95b73c177f276929f03180b02dc77fec38720d718f8361b53d2b8",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							getID(types.BNRoleSyncCommittee),
						),
						nil,
					),
				},
				PostDutyRunnerStateRoot: "ef4080c54e4d1cad83e02388ce61ce499c93628a5176d5defb12c97856dcb4ee",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           errStr,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					testingutils.SSVMsgAggregator(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							getID(types.BNRoleAggregator),
						),
						nil,
					),
				},
				PostDutyRunnerStateRoot: "96cb700dbfff1b478bb550f6b76e66e9837b05e78b77692ca9acc4bebb1845d7",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							getID(types.BNRoleProposer),
						),
						nil,
					),
				},
				PostDutyRunnerStateRoot: "d239eaf711c59e7bbf8c9bf506c95a95051bd528f6bac57cf8465ee1cee821b9",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							getID(types.BNRoleProposer),
						),
						nil,
					),
				},
				PostDutyRunnerStateRoot: "b3e4ccde5539ea5df610151039c16dc5333f8073e1bf658ec8904fc711168219",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
				},
				ExpectedError: errStr,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
							[]types.OperatorID{1, 2, 3},
							2,
							getID(types.BNRoleAttester),
						),
						nil,
					),
				},
				PostDutyRunnerStateRoot: "fdc2d5c52ee8e5dc1a33219579aa820b2e8a22d75bf159cbd68f6299af546e26",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           errStr,
			},
		},
	}
}
