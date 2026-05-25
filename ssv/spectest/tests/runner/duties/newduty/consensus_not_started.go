package newduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ConsensusNotStarted tests starting duty after prev already started but for some duties' consensus didn't start because pre-consensus didnt get quorum (different duties will enable starting a new duty)
func ConsensusNotStarted() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	startRunner := func(r ssv.Runner, duty types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		return r
	}

	multiSpecTest := NewMultiStartNewRunnerDutySpecTest(
		"new duty consensus not started",
		testdoc.NewDutyConsensusNotStartedDoc,
		[]*StartNewRunnerDutySpecTest{
			{
				Name:      "sync committee aggregator",
				Runner:    startRunner(testingutils.AggregatorCommitteeRunner(ks), testingutils.TestingSyncCommitteeContributionDuty),
				Duty:      testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:      "proposer",
				Runner:    startRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionDeneb)),
				Duty:      testingutils.TestingProposerDutyNextEpochV(spec.DataVersionDeneb),
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:      "voluntary exit",
				Runner:    startRunner(testingutils.VoluntaryExitRunner(ks), &testingutils.TestingVoluntaryExitDuty),
				Duty:      &testingutils.TestingVoluntaryExitDutyNextEpoch,
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:      "validator registration",
				Runner:    startRunner(testingutils.ValidatorRegistrationRunner(ks), &testingutils.TestingValidatorRegistrationDuty),
				Duty:      &testingutils.TestingValidatorRegistrationDutyNextEpoch,
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusValidatorRegistrationNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
		},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &StartNewRunnerDutySpecTest{
			Name:      fmt.Sprintf("aggregator (%s)", version.String()),
			Runner:    startRunner(testingutils.AggregatorCommitteeRunner(ks), testingutils.TestingAggregatorDuty(version)),
			Duty:      testingutils.TestingAggregatorDutyNextEpoch(version),
			Threshold: ks.Threshold,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1, version),
				// broadcasts when starting a new duty
			},
		})
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*StartNewRunnerDutySpecTest{
			{
				Name:      fmt.Sprintf("attester (%s)", version.String()),
				Runner:    startRunner(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterDuty(version)),
				Duty:      testingutils.TestingAttesterDutyNextEpoch(version),
				Threshold: ks.Threshold,
			},
			{
				Name:      fmt.Sprintf("sync committee (%s)", version.String()),
				Runner:    startRunner(testingutils.CommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty(version)),
				Duty:      testingutils.TestingSyncCommitteeDutyNextEpoch(version),
				Threshold: ks.Threshold,
			},
			{
				Name:      fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner:    startRunner(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterAndSyncCommitteeDuties(version)),
				Duty:      testingutils.TestingAttesterAndSyncCommitteeDutiesNextEpoch(version),
				Threshold: ks.Threshold,
			},
		}...)
	}

	return multiSpecTest
}
