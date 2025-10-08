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
				Name: "sync committee contribution",
				Runner: decideRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty,
					testingutils.TestSyncCommitteeContributionConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, invalidateSlot(testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)))),
				},
				PostDutyRunnerStateRoot: "f58387d4d4051a2de786e4cbf9dc370a8b19a544f52af04f71195feb3863fc5c",
				DontStartDuty:           true,
				ExpectedErrorCode:       types.FuturePartialSigMsgSlot,
			},
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
				PostDutyRunnerStateRoot: "ff213af6f0bf2350bb37f48021c137dd5552b1c25cb5c6ebd0c1d27debf6080e",
				DontStartDuty:           true,
				ExpectedErrorCode:       types.FuturePartialSigMsgSlot,
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
				PostDutyRunnerStateRoot: "9b4524d5100835df4d71d0a1e559acdc33d541c44a746ebda115c5e7f3eaa85a",
				DontStartDuty:           true,
				ExpectedErrorCode:       types.FuturePartialSigMsgSlot,
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
				PostDutyRunnerStateRoot: "ec573732e70b70808972c43acb5ead6443cff06ba30d8abb51e37ac82ffe0727",
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
				PostDutyRunnerStateRoot: "ec573732e70b70808972c43acb5ead6443cff06ba30d8abb51e37ac82ffe0727",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedVoluntaryExit(ks)),
				},
				ExpectedErrorCode: types.ValidatorExitNoPostConsensusPhaseErrorCode,
			},
			{
				Name: fmt.Sprintf("aggregator (%s)", spec.DataVersionPhase0.String()),
				Runner: decideRunner(
					testingutils.AggregatorRunner(ks),
					testingutils.TestingAggregatorDuty(spec.DataVersionPhase0),
					testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0),
				),
				Duty: testingutils.TestingAggregatorDuty(spec.DataVersionPhase0),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, invalidateSlot(testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1, spec.DataVersionPhase0)))),
				},
				PostDutyRunnerStateRoot: "1fb182fb19e446d61873abebc0ac85a3a9637b51d139cdbd7d8cb70cf7ffec82",
				DontStartDuty:           true,
				ExpectedErrorCode:       types.FuturePartialSigMsgSlot,
			},
			{
				Name: fmt.Sprintf("aggregator (%s)", spec.DataVersionElectra.String()),
				Runner: decideRunner(
					testingutils.AggregatorRunner(ks),
					testingutils.TestingAggregatorDuty(spec.DataVersionElectra),
					testingutils.TestAggregatorConsensusData(spec.DataVersionElectra),
				),
				Duty: testingutils.TestingAggregatorDuty(spec.DataVersionElectra),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, invalidateSlot(testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1, spec.DataVersionElectra)))),
				},
				PostDutyRunnerStateRoot: "1fb182fb19e446d61873abebc0ac85a3a9637b51d139cdbd7d8cb70cf7ffec82",
				DontStartDuty:           true,
				ExpectedErrorCode:       types.InvalidPartialSigMsgSlot,
			},
		},
		ks,
	)

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
