package pre_consensus_justifications

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidNoRunningDuty tests a valid pre-consensus justification for a runner that has no running duty
func ValidNoRunningDuty() *tests.MultiMsgProcessingSpecTest {
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
		Name: "pre consensus justification valid no running duty",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID), nil),
				},
				PostDutyRunnerStateRoot: "17f18a6d0985ea944251bf9b6cc28c2ed5ceb685564fe5598cad8c2654c9b543",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(msgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID), nil),
				},
				PostDutyRunnerStateRoot: "c502e1b533be473b5a79593478dff209751509b7f4907109fbf8852938185db8",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(msgF(testingutils.TestProposerWithJustificationsConsensusData(ks), testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "5ba69a8b0fa59b3afbc64dc85edc32e45347169fb420378a801c6ed8486340fd",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "randao (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(msgF(testingutils.TestProposerBlindedWithJustificationsConsensusData(ks), testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "fb7158179913484f391c8cc2a5efc79e1e195efd99c949aba8f68173965f2b67",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{

				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(msgF(testingutils.TestAttesterConsensusData, testingutils.AttesterMsgID), nil),
				},
				PostDutyRunnerStateRoot: "22297b32fbb3bdabdc5f1b6df5c6e2dfc241ad55f9069f00626b12ebae589912",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(msgF(testingutils.TestSyncCommitteeConsensusData, testingutils.SyncCommitteeMsgID), nil),
				},
				PostDutyRunnerStateRoot: "c67d65112c567bd10a8305f14485f4fcd7eab5258478a75b848be5d8379bfa0b",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
		},
	}
}
