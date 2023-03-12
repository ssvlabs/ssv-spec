package consensus

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoRunningConsensusInstance tests a valid proposal msg before consensus instance starts
func NoRunningConsensusInstance() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedErr := "failed processing consensus message: instance not found"
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus no running consensus instance",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[4], types.OperatorID(4), testingutils.SyncCommitteeContributionMsgID,
							testingutils.TestSyncCommitteeContributionConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "29862cc6054edc8547efcb5ae753290971d664b9c39768503b4d66e1b52ecb06",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.SyncCommitteeMsgID,
							testingutils.TestSyncCommitteeConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "941338be935cec87018b8c076d81ebd28a54e46097352fd4585a89a2bf59b3a2",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.AggregatorMsgID,
							testingutils.TestAggregatorConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "c54e71de23c3957b73abbb0e7b9e195b3f8f6370d62fbec256224faecf177fee",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[4], types.OperatorID(4), testingutils.ProposerMsgID,
							testingutils.TestProposerConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "56eafcb33392ded888a0fefe30ba49e52aa00ab36841cb10c9dc1aa2935af347",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[4], types.OperatorID(4), testingutils.ProposerMsgID,
							testingutils.TestProposerBlindedBlockConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "2ce3241658f324f352c77909f4043934eedf38e939ae638c5ce6acf28e965646",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.AttesterMsgID,
							testingutils.TestAttesterConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "0d5b671f94eeddcb00025dd70fa52d259cafaa5f284645db4fd20e943e2e900d",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
		},
	}
}
