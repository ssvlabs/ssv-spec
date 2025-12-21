package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidMessageSlot tests a valid SignedPartialSignatureMessage with an invalid msg slot
func InvalidMessageSlot() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	invalidateSlot := func(msg *types.PartialSignatureMessages) *types.PartialSignatureMessages {
		msg.Slot = testingutils.TestingDutySlot2
		return msg
	}

	invalidateSlotV := func(msg *types.PartialSignatureMessages, version spec.DataVersion) *types.PartialSignatureMessages {
		msg.Slot = testingutils.TestingInvalidDutySlotV(version)
		return msg
	}

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"post consensus invalid msg slot",
		testdoc.PostConsensusInvalidMsgSlotDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name: "proposer",
				Runner: decideRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
					testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, invalidateSlotV(testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), spec.DataVersionDeneb))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: types.PartialSigMessageFutureSlotErrorCode,
			},
			{
				Name: "proposer (blinded block)",
				Runner: decideRunner(
					testingutils.ProposerBlindedBlockRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
					testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, invalidateSlotV(testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), spec.DataVersionDeneb))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: types.PartialSigMessageFutureSlotErrorCode,
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, invalidateSlot(testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, spec.DataVersionPhase0)))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedValidatorRegistration(ks)),
				},
				ExpectedErrorCode: types.ValidatorRegistrationNoPostConsensusPhaseErrorCode,
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[2], 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[3], 3))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, invalidateSlot(testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, spec.DataVersionPhase0)))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedVoluntaryExit(ks)),
				},
				ExpectedErrorCode: types.ValidatorExitNoPostConsensusPhaseErrorCode,
			},
		},
		ks,
	)

	//// Aggregator committee duty
	//sccSlot := testingutils.TestingSyncCommitteeContributionDuty.Slot
	//multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
	//	Name: "sync committee contribution",
	//	Runner: decideAggregatorCommitteeRunner(
	//		testingutils.AggregatorCommitteeRunner(ks),
	//		testingutils.TestingSyncCommitteeContributionDuty,
	//		testingutils.TestSyncCommitteeContributionConsensusData,
	//	),
	//	Duty: testingutils.TestingSyncCommitteeContributionDuty,
	//	Messages: []*types.SignedSSVMessage{
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, invalidateSlot(testingutils.PostConsensusSyncCommitteeContributionMsgWithSlot(ks.Shares[1], 1, ks, sccSlot)))),
	//	},
	//	DontStartDuty:     true,
	//	ExpectedErrorCode: types.NoRunnerForSlotErrorCode,
	//})
	//
	//for _, version := range testingutils.SupportedAggregatorVersions {
	//	multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
	//		{
	//			Name: fmt.Sprintf("aggregator (%s)", version.String()),
	//			Runner: decideAggregatorCommitteeRunner(
	//				testingutils.AggregatorCommitteeRunner(ks),
	//				testingutils.TestingAggregatorDuty(version),
	//				testingutils.TestAggregatorConsensusData(version),
	//			),
	//			Duty: testingutils.TestingAggregatorDuty(version),
	//			Messages: []*types.SignedSSVMessage{
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, invalidateSlot(testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1, version)))),
	//			},
	//			DontStartDuty:     true,
	//			ExpectedErrorCode: types.NoRunnerForSlotErrorCode,
	//		},
	//		{
	//			Name: fmt.Sprintf("aggregator (%s)", version.String()),
	//			Runner: decideAggregatorCommitteeRunner(
	//				testingutils.AggregatorCommitteeRunner(ks),
	//				testingutils.TestingAggregatorCommitteeDutyMixed(version),
	//				testingutils.TestAggregatorCommitteeConsensusData(version),
	//			),
	//			Duty: testingutils.TestingAggregatorCommitteeDutyMixed(version),
	//			Messages: []*types.SignedSSVMessage{
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, invalidateSlot(testingutils.PostConsensusAggregatorCommitteeMixedMsg(ks.Shares[1], 1, version, ks)))),
	//			},
	//			DontStartDuty:     true,
	//			ExpectedErrorCode: types.NoRunnerForSlotErrorCode,
	//		},
	//	}...)
	//}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{

			{
				Name: fmt.Sprintf("attester (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, invalidateSlot(testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, version)))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: types.NoRunnerForSlotErrorCode,
			},
			{
				Name: fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingSyncCommitteeDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingSyncCommitteeDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, invalidateSlot(testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1, version)))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: types.NoRunnerForSlotErrorCode,
			},
			{
				Name: fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterAndSyncCommitteeDuties(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, invalidateSlot(testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, version)))),
				},
				DontStartDuty:     true,
				ExpectedErrorCode: types.NoRunnerForSlotErrorCode,
			},
		}...)
	}

	return multiSpecTest

}
