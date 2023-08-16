package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidFirstHeight tests a special case for first height which didn't start the instance for duties that require
// preconsesus
func ValidFirstHeight() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msgF := func(obj *types.ConsensusData, id []byte) *qbft.SignedMessage {
		fullData, _ := obj.Encode()
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.Shares[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	firstHeightDuty := func(duty types.Duty) *types.Duty {
		duty.Slot = 0
		return &duty
	}

	firstSlotCd := func(cd *types.ConsensusData) *types.ConsensusData {
		cd.Duty.Slot = 0
		return cd
	}

	//expectedError := "failed processing consensus message: future msg from height, could not process"

	syncContributionDuty := firstHeightDuty(testingutils.TestingSyncCommitteeContributionDuty)
	aggregatorDuty := firstHeightDuty(testingutils.TestingAggregatorDuty)
	proposerDuty := firstHeightDuty(*testingutils.TestingProposerDutyV(spec.DataVersionBellatrix))

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus first height not started",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   syncContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(msgF(
						firstSlotCd(testingutils.TestContributionProofWithJustificationsConsensusDataCustomDuty(ks,
							syncContributionDuty)),
						testingutils.SyncCommitteeContributionMsgID), nil),
				},
				PostDutyRunnerStateRoot: "1c71e2154b44d39541e65319701ce50b3e2d39fab616d2eb06fd72a67bff5793",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   aggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(msgF(
						firstSlotCd(testingutils.TestSelectionProofWithJustificationsConsensusDataCustomDuty(ks, aggregatorDuty)),
						testingutils.AggregatorMsgID), nil),
				},
				PostDutyRunnerStateRoot: "3b9877067deef7be6916fd4879878e51b5047a39d57a804a69589e113c4a893a",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   proposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(msgF(
						firstSlotCd(testingutils.TestProposerWithJustificationsConsensusDataCustomDutyV(ks, proposerDuty,
							spec.DataVersionBellatrix)),
						testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "f7a63280b5e1ccfd430fd6ab9eaf4c7f0bf50b1b03f8d6c2dfdcfe89471d072a",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   firstHeightDuty(*testingutils.TestingProposerDutyV(spec.DataVersionBellatrix)),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(msgF(
						firstSlotCd(testingutils.TestProposerBlindedWithJustificationsConsensusDataCustomDutyV(ks, proposerDuty,
							spec.DataVersionBellatrix)),
						testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "ef8bcdcf507151f25f8247b408e0ad47730c298b62068fff97b3fa8e3b6076c3",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
			},
		},
	}
}
