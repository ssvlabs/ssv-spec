package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Finished tests a valid start duty after finished prev
func Finished() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	finishRunner := func(r ssv.Runner, duty *types.Duty, finishController bool) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)

		// for duties with a consensus controller
		if finishController {
			r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
				r.GetBaseRunner().QBFTController.GetConfig(),
				r.GetBaseRunner().Share,
				r.GetBaseRunner().QBFTController.Identifier,
				qbft.Height(duty.Slot))
			r.GetBaseRunner().State.RunningInstance.State.Decided = true
			r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
			r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.Slot)
		}

		r.GetBaseRunner().State.Finished = true
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty finished",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name: "sync committee aggregator",
				Runner: finishRunner(testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty, true),
				Duty:                    &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "fbcf8689caa2138c117d9385daa53a0c8aae76b9f12014d13a66d984f6287edf",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  finishRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty, true),
				Duty:                    &testingutils.TestingSyncCommitteeDutyNextEpoch,
				PostDutyRunnerStateRoot: "b160f3edf371d390ac09b99f01ae283ffc6339649ca0352f3bd974ab2d6458b3",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  finishRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty, true),
				Duty:                    &testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "e5ffdb2d4b64133979b73370c757b94e2f82952238d7c2bcdc91fbe1782ec80a",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name: "proposer",
				Runner: finishRunner(testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionBellatrix), true),
				Duty:                    testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix),
				PostDutyRunnerStateRoot: "78a040ec69b6932809c8b582999277d3e9766d564eacde7eeb029f0749460f39",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  finishRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty, true),
				Duty:                    &testingutils.TestingAttesterDutyNextEpoch,
				PostDutyRunnerStateRoot: "cbfb9b6302ff1e7a1bf356f57a8e88dd4c4f7ddef6345c62dac125af1d1db4ce",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "voluntary exit",
				Runner:                  finishRunner(testingutils.VoluntaryExitRunner(ks), &testingutils.TestingVoluntaryExitDuty, false),
				Duty:                    &testingutils.TestingVoluntaryExitDutyNextEpoch,
				PostDutyRunnerStateRoot: "6f6d918e15ebc7b84cb77e2d603019d1cbfb6d7293daddd48780da47c14e53ce",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusVoluntaryExitNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "validator registration",
				Runner:                  finishRunner(testingutils.ValidatorRegistrationRunner(ks), &testingutils.TestingValidatorRegistrationDuty, false),
				Duty:                    &testingutils.TestingValidatorRegistrationDutyNextEpoch,
				PostDutyRunnerStateRoot: "6f6d918e15ebc7b84cb77e2d603019d1cbfb6d7293daddd48780da47c14e53ce",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
		},
	}
}
