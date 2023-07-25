package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateDutyFinished is a test that runs the following scenario:
//   - Runner is assigned a duty
//   - Runner finishes the duty
//   - Runner is assigned the same duty again
//     https://github.com/bloxapp/ssv-spec/issues/282.
//     Does it make sense that the runner starts? Does it make sense that duties with precon phase don't return an error?
func DuplicateDutyFinished() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	finishRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.Slot))
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.Slot)
		r.GetBaseRunner().State.Finished = true
		return r
	}

	expectedError := "can't start new duty runner instance for duty: could not start new QBFT instance: instance already running"

	// finishTaskRunner is a helper function that finishes a task runner and returns it
	// task is an operation that isn't a beacon duty, e.g. validator registration
	finishTaskRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.Finished = true
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "duplicate duty finished",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  finishRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionNexEpochDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "a91a90cd4033996ef7ccfccd36dcc2306e56c0b42ab51e6bb97fd88e7d3b6e35",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  finishRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				PostDutyRunnerStateRoot: "07b3b82b83ed99a1bb474dacdb65d26866e1ed9aac13f805849e76133bdb0e2d",
				ExpectedError:           expectedError,
			},
			{
				Name:                    "aggregator",
				Runner:                  finishRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDutyNextEpoch),
				Duty:                    &testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "f58c8a5d953658339cbac025bb1f3d563eaad8693c9af8c780c95d808732e464",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "proposer",
				Runner:                  finishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionBellatrix)),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				PostDutyRunnerStateRoot: "a91e014950037e5dc2ab9e801d0170b90b82f592029a2409c2332f252368d71d",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "attester",
				Runner:                  finishRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "a96148ae850dd3d3a0d63869a95702174739151fa271ba463a3c163cabe35e13",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
			{
				Name: "validator registration",
				Runner: finishTaskRunner(testingutils.ValidatorRegistrationRunner(ks),
					&testingutils.TestingValidatorRegistrationDuty),
				Duty:                    &testingutils.TestingValidatorRegistrationDuty,
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
		},
	}
}
