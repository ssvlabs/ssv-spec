package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests a valid commit msg after returned decided already
func PostDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution),
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingCommitMessageWithIdentifierAndFullData(
							ks.Shares[4], types.OperatorID(4), testingutils.SyncCommitteeContributionMsgID,
							testingutils.TestSyncCommitteeContributionConsensusDataByts,
						), nil)),
				PostDutyRunnerStateRoot: postDecidedSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     postDecidedSyncCommitteeContributionSC().ExpectedState,
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
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee),
					testingutils.SSVMsgSyncCommittee(
						testingutils.TestingCommitMessageWithIdentifierAndFullData(
							ks.Shares[4], types.OperatorID(4), testingutils.SyncCommitteeMsgID,
							testingutils.TestSyncCommitteeConsensusDataByts,
						), nil)),
				PostDutyRunnerStateRoot: postDecidedSyncCommitteeSC().Root(),
				PostDutyRunnerState:     postDecidedSyncCommitteeSC().ExpectedState,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData, ks, types.BNRoleAggregator),
					testingutils.SSVMsgAggregator(
						testingutils.TestingCommitMessageWithIdentifierAndFullData(
							ks.Shares[4], types.OperatorID(4), testingutils.AggregatorMsgID,
							testingutils.TestAggregatorConsensusDataByts,
						), nil)),
				PostDutyRunnerStateRoot: postDecidedAggregatorSC().Root(),
				PostDutyRunnerState:     postDecidedAggregatorSC().ExpectedState,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(spec.DataVersionBellatrix), ks, types.BNRoleProposer),
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMessageWithIdentifierAndFullData(
							ks.Shares[4], types.OperatorID(4), testingutils.ProposerMsgID,
							testingutils.TestProposerConsensusDataBytsV(spec.DataVersionBellatrix),
						), nil)),
				PostDutyRunnerStateRoot: "7bafe77f6aa303e2cd38f741ebd366d68b4ced79ffbd224b98904bb22a58d010",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
				},
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionBellatrix), ks, types.BNRoleProposer),
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMessageWithIdentifierAndFullData(
							ks.Shares[4], types.OperatorID(4), testingutils.ProposerMsgID,
							testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionBellatrix),
						), nil)),
				PostDutyRunnerStateRoot: "d2d3d97a0bb878594c5b5319cbf1a1a6afe2467ef26fb2e47eb6d075b3c90428",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester),
					testingutils.SSVMsgAttester(
						testingutils.TestingCommitMessageWithIdentifierAndFullData(
							ks.Shares[4], types.OperatorID(4), testingutils.AttesterMsgID,
							testingutils.TestAttesterConsensusDataByts,
						), nil)),
				PostDutyRunnerStateRoot: postDecidedAttesterSC().Root(),
				PostDutyRunnerState:     postDecidedAttesterSC().ExpectedState,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
