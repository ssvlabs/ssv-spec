package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PastMessage tests a valid proposal past msg
func PastMessage() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	bumpHeight := func(r ssv.Runner) ssv.Runner {
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().QBFTController.Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight))

		r.GetBaseRunner().QBFTController.Height = 10
		return r
	}

	pastMsgF := func(obj *types.ConsensusData, id []byte) *qbft.SignedMessage {
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

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus past message",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: bumpHeight(testingutils.SyncCommitteeContributionRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(
						pastMsgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "d1ba71cab348c80ebb7b4533c9c482eaba407f6a73864ee742aab93e73b94dab",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "sync committee",
				Runner: bumpHeight(testingutils.SyncCommitteeRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						pastMsgF(testingutils.TestSyncCommitteeConsensusData, testingutils.SyncCommitteeMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "c1fe93bc3530e14e014b78102b078bbdd6e55c2c0979fce40fa00b1400d9e774",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "aggregator",
				Runner: bumpHeight(testingutils.AggregatorRunner(ks)),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(
						pastMsgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "5a1a9b9fb21682ea854f919be531a692fe5c3a6c5302214dbf3421faed57cff8",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "proposer",
				Runner: bumpHeight(testingutils.ProposerRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						pastMsgF(testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionBellatrix), testingutils.ProposerMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "1c939726a237c02013fab61901e819e34ec99e2ef62dadb6c847e5ad118fc4e7",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: bumpHeight(testingutils.ProposerBlindedBlockRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						pastMsgF(testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionBellatrix), testingutils.ProposerMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "49edaab0d759ba8a35a37ab26416ae04962d77ec088b87c4f1e65f781c1ed96f",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "attester",
				Runner: bumpHeight(testingutils.AttesterRunner(ks)),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						pastMsgF(testingutils.TestAttesterConsensusData, testingutils.AttesterMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "7218b8c7b0c3fc4751b62ed15960de36bb28d6ec800a6d689d16ae99754fbf66",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.ValidatorRegistrationMsgID,
							testingutils.TestAttesterConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "no consensus phase for validator registration",
			},
		},
	}
}
