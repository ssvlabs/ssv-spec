package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidNoRunningDuty tests a valid pre-consensus justification for a runner that has no running duty
func ValidNoRunningDuty() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msgF := func(obj *types.ConsensusData, id []byte) *types.SignedSSVMessage {
		fullData, _ := obj.Encode()
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     1,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.OperatorKeys[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus justification valid no running duty",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: decideFirstHeight(testingutils.SyncCommitteeContributionRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID),
				},
				PostDutyRunnerStateRoot: "118bb07beaee54ffa3f4a676f752cdf391692b5981066b4827a8c29d100756ed",
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
			},
			{
				Name:   "aggregator selection proof",
				Runner: decideFirstHeight(testingutils.AggregatorRunner(ks)),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID),
				},
				PostDutyRunnerStateRoot: "8739b666e50accff30ff3f3b6fe6f75eba4d0eec340efd008dc9f66239155292",
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
			},
			{
				Name:   "randao",
				Runner: decideFirstHeight(testingutils.ProposerRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb), testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "a7762b4ffea3d8fc26474348d8222a8ab13c985bbe5a1140b115ce1788acf7e3",
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
			},
			{
				Name:   "randao (blinded block)",
				Runner: decideFirstHeight(testingutils.ProposerBlindedBlockRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb), testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "1d3ddfd3bce80e7795ed4bfe3a23b943f395dcdb0f6d0844f25a149e6d1bea28",
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
			},
			{

				Name:   "attester",
				Runner: decideFirstHeight(testingutils.CommitteeRunner(ks)),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestAttesterConsensusData, testingutils.AttesterMsgID),
				},
				PostDutyRunnerStateRoot: "97ba097765d1b3e9e1b0dde97723d5ec3100e12ed197d567cdba32315409d20b",
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
			},
			{
				Name:   "sync committee",
				Runner: decideFirstHeight(testingutils.SyncCommitteeRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestSyncCommitteeConsensusData, testingutils.SyncCommitteeMsgID),
				},
				PostDutyRunnerStateRoot: "631bde3df3e21a17dac94c01623d58cc02b1c99f52120e85f0592c3d2626dddd",
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
			},
		},
	}
}
