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
				PostDutyRunnerStateRoot: "fad43e18e8df967df0cffc06f4947951ccce347704db52577cbdfb17a6d87f50",
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
				PostDutyRunnerStateRoot: "664957056215c863ed63cc9879bab9358b210a3fd96032e51a6cf16c6f4deeab",
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
				PostDutyRunnerStateRoot: "cbf71df66ebbaa94209a0f0d8ddfc326bf10395c331e238ad8aba20b5c0762ef",
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
				PostDutyRunnerStateRoot: "89ad4f86c8280b8538d8048807ac60ebe8fa4f830275da2910c1862c3d4c3567",
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
				PostDutyRunnerStateRoot: "3e9c594d4bd4d0c1afe681b0e4698cd645a76e5ead41263b3cfa0ffa2ab34754",
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
				PostDutyRunnerStateRoot: "bce6c4862a058b963f30a5c2eebefcd244bc156cc21a8bf84970cc38f703673d",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
