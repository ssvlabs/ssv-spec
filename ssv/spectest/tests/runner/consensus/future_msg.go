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
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					futureMsgF(testingutils.TestSyncCommitteeContributionConsensusData, testingutils.SyncCommitteeContributionMsgID),
				},
				PostDutyRunnerStateRoot: "68fd25b1cb30902e7b7b3e7ff674c3862ff956954a06fac0df485961b8bb3934",
				DontStartDuty:           true,
				ExpectedErrorCode:       expectedErrorCode,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					futureMsgF(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb), testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "32dd1d1d7a4c34bb7dafc0866f69eb49f6a0a23755b135f83ad14d12e39fff82",
				DontStartDuty:           true,
				ExpectedErrorCode:       expectedErrorCode,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					futureMsgF(testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb),
						testingutils.ProposerMsgID),
				},
				PostDutyRunnerStateRoot: "58b946451dc5ccbd52fbc9e6bbe0ac888253d1708be018a3ff0b07762dd28891",
				DontStartDuty:           true,
				ExpectedErrorCode:       expectedErrorCode,
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
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
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
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
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
			Runner: testingutils.AggregatorRunner(ks),
			Duty:   testingutils.TestingAggregatorDuty(version),
			Messages: []*types.SignedSSVMessage{
				futureMsgF(testingutils.TestAggregatorConsensusData(version), testingutils.AggregatorMsgID),
			},
			PostDutyRunnerStateRoot: "bdc7c2150e0f2d4669e112848f5140b52aba0367b60ff2b594d5a5bef3587834",
			DontStartDuty:           true,
			ExpectedErrorCode:       expectedErrorCode,
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
