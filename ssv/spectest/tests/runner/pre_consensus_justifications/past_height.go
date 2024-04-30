package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PastHeight tests justification with height <= current height (not valid)
func PastHeight() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msgF := func(obj *types.ConsensusData, id []byte) *types.SignedSSVMessage {
		fullData, _ := obj.Encode()
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.OperatorKeys[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus past height",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: decideFirstHeight(testingutils.SyncCommitteeContributionRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID),
				},
				PostDutyRunnerStateRoot: "4779f0f8a875748eda9fd96ba1fdb4e9cec53211c17c8898af47f36436b89833",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "aggregator selection proof",
				Runner: decideFirstHeight(testingutils.AggregatorRunner(ks)),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID),
				},
				PostDutyRunnerStateRoot: "72188c3cc54996765881db396e36d32090053001b4922f54cfd85d6748c4ce6a",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao",
				Runner: decideFirstHeight(testingutils.ProposerRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb), testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "c41cabb4acc27875342979ee8e77be17b29eea25442fe97c8c3a16711329c03f",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao (blinded block)",
				Runner: decideFirstHeight(testingutils.ProposerBlindedBlockRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb), testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "582b8769e018c48b9263655cc315d3b889e921d0bef59537325c414c1ca3bace",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
			},
			// {

			// 	Name:   "attester",
			// 	Runner: decideFirstHeight(testingutils.CommitteeRunner(ks)),
			// 	Duty:   &testingutils.TestingAttesterDuty,
			// 	Messages: []*types.SignedSSVMessage{
			// 		msgF(testingutils.TestAttesterConsensusData, testingutils.AttesterMsgID),
			// 	},
			// 	PostDutyRunnerStateRoot: "a60c5cb9e8c0189086ae447ee5d1f5477a0be3eb3af4fb15f17679e930f52538",
			// 	OutputMessages:          []*types.PartialSignatureMessages{},
			// },
			// {
			// 	Name:   "sync committee",
			// 	Runner: decideFirstHeight(testingutils.SyncCommitteeRunner(ks)),
			// 	Duty:   &testingutils.TestingSyncCommitteeDuty,
			// 	Messages: []*types.SignedSSVMessage{
			// 		msgF(testingutils.TestSyncCommitteeConsensusData, testingutils.SyncCommitteeMsgID),
			// 	},
			// 	PostDutyRunnerStateRoot: "81118c674d65bddd1c468e4ec6bb34d2d9355e2008083f3fed315e7af9061fd5",
			// 	OutputMessages:          []*types.PartialSignatureMessages{},
			// },
		},
	}
}
