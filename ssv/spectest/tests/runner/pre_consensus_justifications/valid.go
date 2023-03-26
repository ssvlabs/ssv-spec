package pre_consensus_justifications

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Valid tests a valid pre-consensus justification
func Valid() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	msgF := func(obj *types.ConsensusData, id []byte) *qbft.SignedMessage {
		fullData, _ := obj.Encode()
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     1,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.Shares[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus valid quorum",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID), nil),
				},
				PostDutyRunnerStateRoot: "f1cdb47ee453decf08c04161f463d9001a5eb197e348486190cebff438c75e35",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(msgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID), nil),
				},
				PostDutyRunnerStateRoot: "8d811303faae71ad17666fef82e0692c695ab9388d436385be46a5821271f409",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(msgF(testingutils.TestProposerWithJustificationsConsensusData(ks), testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "aa3cf0b43cea31e9c4bc13b5e4cbb150b362d15fffd5c8d5fdc4848bb4eff638",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(msgF(testingutils.TestProposerBlindedWithJustificationsConsensusData(ks), testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "81d34e0b6c9a61243b607b194527313579e07499abafd82ed0a34dfecdab408e",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{

				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(msgF(testingutils.TestAttesterConsensusData, testingutils.AttesterMsgID), nil),
				},
				PostDutyRunnerStateRoot: "757b432ae782aca13549cdaeb5df8b292691bd8eaab0057ae779c59ed222fd79",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(msgF(testingutils.TestSyncCommitteeConsensusData, testingutils.SyncCommitteeMsgID), nil),
				},
				PostDutyRunnerStateRoot: "fa35779ebcdb59b635aa8e661cbbf006a7bae1eded677287001277a4205631a2",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
		},
	}
}
