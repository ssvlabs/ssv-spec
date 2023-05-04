package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// PostFinish tests a valid commit msg after runner finished
func PostFinish() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post finish",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution),
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.SyncCommitteeContributionMsgID,
							testingutils.TestSyncCommitteeContributionConsensusDataByts,
						), nil),
				),
				PostDutyRunnerStateRoot: "5448f47bb76e4639629e146c242cb27a4a265fae9d871f0fc0f3c66aeea60eb0",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee),
					testingutils.SSVMsgSyncCommittee(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.SyncCommitteeMsgID,
							testingutils.TestSyncCommitteeConsensusDataByts,
						), nil),
				),
				PostDutyRunnerStateRoot: "8097eefe5cbd4590de980ac66db0f033552f608fef580d1b9c452bcf2f1513fd",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestAggregatorConsensusData, ks, types.BNRoleAggregator),
					testingutils.SSVMsgAggregator(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.AggregatorMsgID,
							testingutils.TestAggregatorConsensusDataByts,
						), nil),
				),
				PostDutyRunnerStateRoot: "99ddd0aadc2708e0a33f9dd979fd45af9bf71d8d85fc1b04cbb0e418e909bcd4",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestProposerConsensusData, ks, types.BNRoleProposer),
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.ProposerMsgID,
							testingutils.TestProposerConsensusDataByts,
						), nil),
				),
				PostDutyRunnerStateRoot: "7bafe77f6aa303e2cd38f741ebd366d68b4ced79ffbd224b98904bb22a58d010",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestProposerBlindedBlockConsensusData, ks, types.BNRoleProposer),
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.ProposerMsgID,
							testingutils.TestProposerBlindedBlockConsensusDataByts,
						), nil),
				),
				PostDutyRunnerStateRoot: "d2d3d97a0bb878594c5b5319cbf1a1a6afe2467ef26fb2e47eb6d075b3c90428",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester),
					testingutils.SSVMsgAttester(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.AttesterMsgID,
							testingutils.TestAttesterConsensusDataByts,
						), nil),
				),
				PostDutyRunnerStateRoot: "4016e9bc1405a443f4a2755f5927d9017c66dd1f5246d03c9bcf7b353649a460",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
