package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidFirstHeight tests a special case for first height which didn't start
func ValidFirstHeight() tests.SpecTest {
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
		Name: "pre consensus first height not started",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID),
				},
				PostDutyRunnerStateRoot: "1c71e2154b44d39541e65319701ce50b3e2d39fab616d2eb06fd72a67bff5793",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID),
				},
				PostDutyRunnerStateRoot: "3b9877067deef7be6916fd4879878e51b5047a39d57a804a69589e113c4a893a",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb), testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "f7a63280b5e1ccfd430fd6ab9eaf4c7f0bf50b1b03f8d6c2dfdcfe89471d072a",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					msgF(testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb), testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "ef8bcdcf507151f25f8247b408e0ad47730c298b62068fff97b3fa8e3b6076c3",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
			},
			// {

			// 	Name:   "attester",
			// 	Runner: testingutils.CommitteeRunner(ks),
			// 	Duty:   &testingutils.TestingAttesterDuty,
			// 	Messages: []*types.SignedSSVMessage{
			// 		msgF(testingutils.TestAttesterConsensusData, testingutils.AttesterMsgID),
			// 	},
			// 	PostDutyRunnerStateRoot: "0d5b671f94eeddcb00025dd70fa52d259cafaa5f284645db4fd20e943e2e900d",
			// 	OutputMessages:          []*types.PartialSignatureMessages{},
			// },
			// {
			// 	Name:   "sync committee",
			// 	Runner: testingutils.SyncCommitteeRunner(ks),
			// 	Duty:   &testingutils.TestingSyncCommitteeDuty,
			// 	Messages: []*types.SignedSSVMessage{
			// 		msgF(testingutils.TestSyncCommitteeConsensusData, testingutils.SyncCommitteeMsgID),
			// 	},
			// 	PostDutyRunnerStateRoot: "5adbf2c86193070a8f74596275e7a62d48a6a573259150d7ec694b3571c7a787",
			// 	OutputMessages:          []*types.PartialSignatureMessages{},
			// },
		},
	}
}
