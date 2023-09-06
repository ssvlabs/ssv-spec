package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DecidedInstance tests a valid pre-consensus justification for a decided instance
// pre-consensus will return false from shouldProcessingJustificationsForHeight as it's decided
//

func DecidedInstance() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msgF := func(obj *types.ConsensusData, id []byte) *qbft.SignedMessage {
		fullData, _ := obj.Encode()
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     testingutils.TestingDutySlot,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.Shares[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus decided instance",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SSVDecideSyncCommiteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID), nil),
				},
				PostDutyRunnerStateRoot: "4779f0f8a875748eda9fd96ba1fdb4e9cec53211c17c8898af47f36436b89833",
				DontStartDuty:           true,
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.SSVDecideAggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(msgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID), nil),
				},
				PostDutyRunnerStateRoot: "72188c3cc54996765881db396e36d32090053001b4922f54cfd85d6748c4ce6a",
				DontStartDuty:           true,
			},
			{
				Name:   "randao",
				Runner: testingutils.SSVDecideProposerRunnerV(ks, spec.DataVersionBellatrix),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(msgF(testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionBellatrix), testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "c41cabb4acc27875342979ee8e77be17b29eea25442fe97c8c3a16711329c03f",
				DontStartDuty:           true,
			},
			{
				Name:   "randao (blinded block)",
				Runner: testingutils.SSVDecideBlindedProposerRunnerV(ks, spec.DataVersionBellatrix),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(msgF(testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionBellatrix), testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "582b8769e018c48b9263655cc315d3b889e921d0bef59537325c414c1ca3bace",
				DontStartDuty:           true,
			},
			{

				Name:   "attester",
				Runner: testingutils.SSVDecideAttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(msgF(testingutils.TestAttesterConsensusData, testingutils.AttesterMsgID), nil),
				},
				PostDutyRunnerStateRoot: "a60c5cb9e8c0189086ae447ee5d1f5477a0be3eb3af4fb15f17679e930f52538",
				DontStartDuty:           true,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SSVDecideSyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(msgF(testingutils.TestSyncCommitteeConsensusData, testingutils.SyncCommitteeMsgID), nil),
				},
				PostDutyRunnerStateRoot: "81118c674d65bddd1c468e4ec6bb34d2d9355e2008083f3fed315e7af9061fd5",
				DontStartDuty:           true,
			},
		},
	}
}
