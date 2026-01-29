package consensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FutureMessage tests a valid proposal future msg
func FutureMessage() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	futureMsgF := func(value types.Encoder, id []byte) *types.SignedSSVMessage {
		var fullData []byte
		if value != nil {
			fullData, _ = value.Encode()
		} else {
			panic("no consensus data or beacon vote")
		}
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     10,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.OperatorKeys[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	expectedErrorCode := types.FutureMessageErrorCode
	expectedErrorCommitteeCode := types.NoRunnerForSlotErrorCode

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"consensus future message",
		testdoc.ConsensusFutureMessageDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.AggregatorCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					futureMsgF(testingutils.TestSyncCommitteeContributionConsensusData, testingutils.TestingAggregatorCommitteeMsgID[:]),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCommitteeCode,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					futureMsgF(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb), testingutils.ProposerMsgID),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					futureMsgF(testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb),
						testingutils.ProposerMsgID),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1),
						testingutils.ValidatorRegistrationMsgID, testingutils.TestAttesterConsensusDataByts,
						qbft.Height(testingutils.TestingDutySlot)),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: types.ValidatorRegistrationNoConsensusPhaseErrorCode,
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1),
						testingutils.VoluntaryExitMsgID, testingutils.TestAttesterConsensusDataByts,
						qbft.Height(testingutils.TestingDutySlot)),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedErrorCode: types.ValidatorExitNoConsensusPhaseErrorCode,
			},
		},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("aggregator (%s)", version.String()),
			Runner: testingutils.AggregatorCommitteeRunner(ks),
			Duty:   testingutils.TestingAggregatorDuty(version),
			Messages: []*types.SignedSSVMessage{
				futureMsgF(testingutils.TestAggregatorConsensusData(version), testingutils.TestingAggregatorCommitteeMsgID[:]),
			},
			DontStartDuty:     true,
			ExpectedErrorCode: expectedErrorCommitteeCode,
		},
		)
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name:   fmt.Sprintf("attester (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterDuty(version),
				Messages: []*types.SignedSSVMessage{
					futureMsgF(&testingutils.TestBeaconVote, testingutils.CommitteeMsgID(ks)),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCommitteeCode,
			},
			{
				Name:   fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty(version),
				Messages: []*types.SignedSSVMessage{
					futureMsgF(&testingutils.TestBeaconVote, testingutils.CommitteeMsgID(ks)),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCommitteeCode,
			},
			{
				Name:   fmt.Sprintf("attester sync committee (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages: []*types.SignedSSVMessage{
					futureMsgF(&testingutils.TestBeaconVote, testingutils.CommitteeMsgID(ks)),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCommitteeCode,
			},
		}...)
	}

	return multiSpecTest
}
