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

// DuplicateDutyNotFinished is a test that runs the following scenario:
// - Runner is assigned a duty
// - Runner doesn't finish the duty
// - Runner is assigned the same duty again
func DuplicateDutyNotFinished() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	notFinishRunner := func(r ssv.Runner, duty types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().QBFTController.CommitteeMember,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.DutySlot()))
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.DutySlot())
		return r
	}

	// notFinishTaskRunner is a helper function that finishes a task runner and returns it
	// task is an operation that isn't a beacon duty, e.g. validator registration
	notFinishTaskRunner := func(r ssv.Runner, duty *types.ValidatorDuty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		return r
	}

	expectedError := fmt.Sprintf("can't start duty: duty for slot %d already passed. Current height is %d",
		testingutils.TestingDutySlot,
		testingutils.TestingDutySlot)

	expectedTaskError := fmt.Sprintf("can't start non-beacon duty: duty for slot %d already passed. "+
		"Current slot is %d",
		testingutils.TestingDutySlot,
		testingutils.TestingDutySlot)

	return &MultiStartNewRunnerDutySpecTest{
		Name: "duplicate duty not finished",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  notFinishRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "f8f6de434622433553716a8d16abf43d92ddad8fd33f0350c39f87d12e30d7e2",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "aggregator",
				Runner:                  notFinishRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "c47ae8ab2b504bcf0439e541f2d9d04138ccb45fee9411e12da9fd857a46a5b6",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "proposer",
				Runner:                  notFinishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionDeneb)),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "914d74606f9de8a2425b875d248532564e1770a6320f923ecad1dd12998b1158",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedError: fmt.Sprintf("can't start duty: duty for slot %d already passed. Current height is %d",
					testingutils.TestingDutySlotV(spec.DataVersionDeneb),
					testingutils.TestingDutySlotV(spec.DataVersionDeneb)),
			},
			{
				Name:                    "attester",
				Runner:                  notFinishRunner(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterDuty),
				Duty:                    testingutils.TestingAttesterDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "c3556c0d6524a6483057916e68c49e8815b25a47cf8e6677c5a37c2a42f89629",
				OutputMessages:          []*types.PartialSignatureMessages{},
				ExpectedError:           expectedError,
			},
			{
				Name:                    "sync committee",
				Runner:                  notFinishRunner(testingutils.CommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "c3556c0d6524a6483057916e68c49e8815b25a47cf8e6677c5a37c2a42f89629",
				OutputMessages:          []*types.PartialSignatureMessages{},
				ExpectedError:           expectedError,
			},
			{
				Name:                    "attester and sync committee",
				Runner:                  notFinishRunner(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterAndSyncCommitteeDuties),
				Duty:                    testingutils.TestingAttesterAndSyncCommitteeDuties,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "c3556c0d6524a6483057916e68c49e8815b25a47cf8e6677c5a37c2a42f89629",
				OutputMessages:          []*types.PartialSignatureMessages{},
				ExpectedError:           expectedError,
			},
			{
				Name: "validator registration",
				Runner: notFinishTaskRunner(testingutils.ValidatorRegistrationRunner(ks),
					&testingutils.TestingValidatorRegistrationDuty),
				Duty:                    &testingutils.TestingValidatorRegistrationDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages:          []*types.PartialSignatureMessages{},
				ExpectedError:           expectedTaskError,
			},
			{
				Name: "voluntary exit",
				Runner: notFinishTaskRunner(testingutils.VoluntaryExitRunner(ks),
					&testingutils.TestingVoluntaryExitDuty),
				Duty:                    &testingutils.TestingVoluntaryExitDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages:          []*types.PartialSignatureMessages{},
				ExpectedError:           expectedTaskError,
			},
		},
	}
}
