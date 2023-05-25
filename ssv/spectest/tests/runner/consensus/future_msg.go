package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FutureMessage tests a valid proposal future msg
func FutureMessage() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	futureMsgF := func(obj *types.ConsensusData, id []byte) *qbft.SignedMessage {
		fullData, _ := obj.Encode()
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     10,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.Shares[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus future message",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(
						futureMsgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "865a12c7c4323f1ed79c19a072ccf9b8d8f37509dc04245ace97fb3b2cf85dc9",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						futureMsgF(testingutils.TestSyncCommitteeConsensusData, testingutils.SyncCommitteeMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "40595793489df2e1e8ca69738b24b3aefdce249d486793d59397c019455a67c5",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(
						futureMsgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "a7b4f13138df18996539187cbfc55e625ea29734de4ccf49cc165d02c1dea3f4",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						futureMsgF(testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionBellatrix), testingutils.ProposerMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "3cc9da8e55d6a3ff6ea106a04d892c02f517c61a803032e8c33b27fd5afd1b7b",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						futureMsgF(testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionBellatrix), testingutils.ProposerMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "7066954dc0e63bbdf180e8d34c97f1bf46e00bf26c1f9f1172238fc7b30ae437",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						futureMsgF(testingutils.TestAttesterConsensusData, testingutils.AttesterMsgID),
						nil),
				},
				PostDutyRunnerStateRoot: "dae20f18cb0b47fac3cf4f4782b633e4871fa63a7c6584cd61860f7d5818abf9",
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
				PostDutyRunnerStateRoot: "72cbd4d0662f728990371054688225487e357fb0de630e2edc8cb13c551b02af",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "no consensus phase for validator registration",
			},
		},
	}
}
