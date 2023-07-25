package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostWrongDecided tests starting a new duty after prev was decided wrongly (future decided)
// This can happen if we receive a future decided message from the network.
func PostWrongDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// https://github.com/bloxapp/ssv-spec/issues/285. We initialize the runner with an impossible decided value.
	// Maybe we should ensure that `ValidateDecided()` doesn't let the runner enter this state and delete the test?
	decideWrong := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		storedInstances := r.GetBaseRunner().QBFTController.StoredInstances
		storedInstances = append(storedInstances, nil)
		storedInstances = append(storedInstances, nil)

		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		storedInstances[1] = r.GetBaseRunner().State.RunningInstance

		higherDecided := qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			50)
		higherDecided.State.Decided = true
		higherDecided.State.DecidedValue = []byte{1, 2, 3, 4}
		storedInstances[0] = higherDecided
		r.GetBaseRunner().QBFTController.Height = 50
		// TODO: hacky fix to a bug in the test.
		// You can't append a copied slice and expect the original to change in go. Since maybe we want to delete
		// the test I didn't do it nicer.
		r.GetBaseRunner().QBFTController.StoredInstances = storedInstances
		return r
	}

	expectedError := "can't start new duty runner instance for duty: could not start new QBFT instance: attempting to start an instance with a past height"

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post wrong decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  decideWrong(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				PostDutyRunnerStateRoot: "1755f6cbd21dd492f8398bc09909d2f5509d11c9018d6582fe4d94c4fd1473a4",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  decideWrong(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "037414235648efdf794077306b597690be309e4baba7cfd687d6e90b1bbe2402",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
			{
				Name:                    "aggregator",
				Runner:                  decideWrong(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "441378c466c16e0b59be45dbaaa7321990a8a8df1bfc8d4cba0b4dac0efe76fc",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  decideWrong(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionBellatrix)),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				PostDutyRunnerStateRoot: "f1f8ef52d758ac2bd1222943dad8c7a3aaf8bfd54255a65f192cc36791b8e862",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  decideWrong(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "43819363c8d1080883c8d4ecb2b25fc1cde06facb14825c533a891f0be6f539a",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
		},
	}
}
