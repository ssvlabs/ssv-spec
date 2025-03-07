package newduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Finished tests a valid start duty after finished prev
func Finished() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	finishRunner := func(r ssv.Runner, duty types.Duty, finishController bool) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)

		// for duties with a consensus controller
		if finishController {
			r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
				r.GetBaseRunner().QBFTController.GetConfig(),
				r.GetBaseRunner().QBFTController.CommitteeMember,
				r.GetBaseRunner().QBFTController.Identifier,
				qbft.Height(duty.DutySlot()),
				r.GetBaseRunner().QBFTController.OperatorSigner)
			r.GetBaseRunner().State.RunningInstance.State.Decided = true
			r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
			r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.DutySlot())
		}

		r.GetBaseRunner().State.Finished = true
		return r
	}

	multiSpecTest := &MultiStartNewRunnerDutySpecTest{
		Name: "new duty finished",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name: "sync committee aggregator",
				Runner: finishRunner(testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty, true),
				Duty:      &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name: "proposer",
				Runner: finishRunner(testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb), true),
				Duty:      testingutils.TestingProposerDutyNextEpochV(spec.DataVersionDeneb),
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
			},
			{
				Name:      "voluntary exit",
				Runner:    finishRunner(testingutils.VoluntaryExitRunner(ks), &testingutils.TestingVoluntaryExitDuty, false),
				Duty:      &testingutils.TestingVoluntaryExitDutyNextEpoch,
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:      "validator registration",
				Runner:    finishRunner(testingutils.ValidatorRegistrationRunner(ks), &testingutils.TestingValidatorRegistrationDuty, false),
				Duty:      &testingutils.TestingValidatorRegistrationDutyNextEpoch,
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusValidatorRegistrationNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
		},
	}

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &StartNewRunnerDutySpecTest{
			Name:      fmt.Sprintf("aggregator (%s)", version.String()),
			Runner:    finishRunner(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDuty(version), true),
			Duty:      testingutils.TestingAggregatorDutyNextEpoch(version),
			Threshold: ks.Threshold,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1, version), // broadcasts when starting a new duty
			},
		},
		)
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*StartNewRunnerDutySpecTest{

			{
				Name:           fmt.Sprintf("attester (%s)", version.String()),
				Runner:         finishRunner(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterDuty(version), true),
				Duty:           testingutils.TestingAttesterDutyNextEpoch(version),
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
			},
			{
				Name:           fmt.Sprintf("sync committee (%s)", version.String()),
				Runner:         finishRunner(testingutils.CommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty(version), true),
				Duty:           testingutils.TestingSyncCommitteeDutyNextEpoch(version),
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
			},
			{
				Name:           fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner:         finishRunner(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterAndSyncCommitteeDuties(version), true),
				Duty:           testingutils.TestingAttesterAndSyncCommitteeDutiesNextEpoch(version),
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		}...)
	}

	return multiSpecTest
}
