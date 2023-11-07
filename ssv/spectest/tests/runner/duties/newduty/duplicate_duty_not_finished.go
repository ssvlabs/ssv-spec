package newduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateDutyNotFinished is a test that runs the following scenario:
// - Runner is assigned a duty
// - Runner doesn't finish the duty
// - Runner is assigned the same duty again
func DuplicateDutyNotFinished() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	notFinishRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.Slot))
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.Slot)
		return r
	}

	// notFinishTaskRunner is a helper function that finishes a task runner and returns it
	// task is an operation that isn't a beacon duty, e.g. validator registration
	notFinishTaskRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
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
				PostDutyRunnerStateRoot: "f8f6de434622433553716a8d16abf43d92ddad8fd33f0350c39f87d12e30d7e2",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "sync committee",
				Runner:                  notFinishRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				PostDutyRunnerStateRoot: "1bc2227b9a53699b42f5581911ef9b1a51ba2fbe481449195739f17fb61b4178",
				ExpectedError:           expectedError,
			},
			{
				Name:                    "aggregator",
				Runner:                  notFinishRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "c47ae8ab2b504bcf0439e541f2d9d04138ccb45fee9411e12da9fd857a46a5b6",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "proposer",
				Runner:                  notFinishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionBellatrix)),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				PostDutyRunnerStateRoot: "914d74606f9de8a2425b875d248532564e1770a6320f923ecad1dd12998b1158",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "attester",
				Runner:                  notFinishRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "c3556c0d6524a6483057916e68c49e8815b25a47cf8e6677c5a37c2a42f89629",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
			{
				Name: "validator registration",
				Runner: notFinishTaskRunner(testingutils.ValidatorRegistrationRunner(ks),
					&testingutils.TestingValidatorRegistrationDuty),
				Duty:                    &testingutils.TestingValidatorRegistrationDuty,
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedTaskError,
			},
			{
				Name: "voluntary exit",
				Runner: notFinishTaskRunner(testingutils.VoluntaryExitRunner(ks),
					&testingutils.TestingVoluntaryExitDuty),
				Duty:                    &testingutils.TestingVoluntaryExitDuty,
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedTaskError,
			},
		},
	}
}
