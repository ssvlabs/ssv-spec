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
func PostWrongDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

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
			10)
		higherDecided.State.Decided = true
		higherDecided.State.DecidedValue = []byte{1, 2, 3, 4}
		storedInstances[0] = higherDecided
		r.GetBaseRunner().QBFTController.Height = 10
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post wrong decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  decideWrong(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				PostDutyRunnerStateRoot: "8edfaf138366f4b5c833a3d788ea1cd908e2689467cf9e1f8088fe63fa18b514",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  decideWrong(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "e87551c3cf6bb00611d4d48ce60676de620f437fbc4c196994bb1dfd4ccc37c0",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  decideWrong(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "6cbb72768f30754b14d9901d33ebe703b40e92701c9b4345c00aad74895a4d6a",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  decideWrong(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionBellatrix)),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				PostDutyRunnerStateRoot: "a5f243eaf3e242dd9b61d7970a8c5761610c2dcb8928a6027628086c3b3ea586",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  decideWrong(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "066c6d6e76b264ea4774281d01067b402132ef3a44a4e0864dc714584fc8e3a0",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
		},
	}
}
