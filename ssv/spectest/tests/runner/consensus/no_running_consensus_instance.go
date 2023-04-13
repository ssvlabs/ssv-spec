package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoRunningConsensusInstance tests a valid proposal msg before consensus instance starts
func NoRunningConsensusInstance() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	// register runners and return roots
	roots := noRunningConsensusInstanceStateComparison().Register().Roots()

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
						testingutils.TestingPrepareMessageWithParams(
							ks.Shares[1],
							1,
							qbft.FirstRound,
							qbft.FirstHeight,
							testingutils.SyncCommitteeContributionMsgID,
							testingutils.TestingQBFTRootData,
						),
						nil),
				},
				PostDutyRunnerStateRoot: roots[0], // "b2e883cc81caaed04f3e40e8561ae55aa1f6abcdb3168e5cc5c834b1d327026e",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedErr,
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
				PostDutyRunnerStateRoot: roots[1], // "5adbf2c86193070a8f74596275e7a62d48a6a573259150d7ec694b3571c7a787",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedErr,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(
						testingutils.TestingPrepareMessageWithParams(
							ks.Shares[1],
							1,
							qbft.FirstRound,
							qbft.FirstHeight,
							testingutils.AggregatorMsgID,
							testingutils.TestingQBFTRootData,
						),
						nil),
				},
				PostDutyRunnerStateRoot: roots[2], // "6e1095601c6fbbd6ba5912dfe296b50db2ae67d4115bce7aa2ad0b091c693ea5",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedErr,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.TestingPrepareMessageWithParams(
							ks.Shares[1],
							1,
							qbft.FirstRound,
							qbft.FirstHeight,
							testingutils.ProposerMsgID,
							testingutils.TestingQBFTRootData,
						),
						nil),
				},
				PostDutyRunnerStateRoot: roots[3], // "a70427708a0ab6995225538b39e7de5cb622af9651fb02a162c6bfbdf5d0966d",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedErr,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.TestingPrepareMessageWithParams(
							ks.Shares[1],
							1,
							qbft.FirstRound,
							qbft.FirstHeight,
							testingutils.ProposerMsgID,
							testingutils.TestingQBFTRootData,
						),
						nil),
				},
				PostDutyRunnerStateRoot: roots[4], // "dc9ee0b1b1d1562763855898c9962957bc5d4f3090890419c22e0162705e9ca0",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedErr,
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
				PostDutyRunnerStateRoot: roots[5], // "0d5b671f94eeddcb00025dd70fa52d259cafaa5f284645db4fd20e943e2e900d",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedErr,
			},
		},
	}
}
