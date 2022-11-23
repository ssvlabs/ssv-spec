package newduty

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests a valid start duty before finished and after decided
func PostDecided() *MultiStartNewRunnerDutySpecTest {
	ks := testingutils.Testing4SharesSet()

	preInstances := func(r ssv.Runner) []*qbft.Instance {
		runningInstance := qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		runningInstance.State.Decided = true
		return []*qbft.Instance{runningInstance}
	}

	// TODO: check error
	// nolint
	finishRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningHeight = qbft.FirstHeight
		r.GetBaseRunner().QBFTController.Height = qbft.FirstHeight

		// support running tests without json
		for _, instance := range preInstances(r) {
			if err := r.GetBaseRunner().QBFTController.SaveInstance(instance); err != nil {
				panic(err)
			}
		}
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  finishRunner(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionNexEpochDuty),
				PreStoredInstances:      preInstances(testingutils.SyncCommitteeContributionRunner(ks)),
				Duty:                    testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "e0efcdf6fa389d2f51402459953c1fec84a0bef58cfbc6fbd2e45fb01d5937eb",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  finishRunner(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				PreStoredInstances:      preInstances(testingutils.SyncCommitteeRunner(ks)),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "51867eaee7af003ec2c2f16853bb3d0fbe3bd7c10ba3bdab7feb12f5245e0cec",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  finishRunner(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDutyNextEpoch),
				PreStoredInstances:      preInstances(testingutils.AggregatorRunner(ks)),
				Duty:                    testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "b5867556db1b759847704be66758c5979e69345c30dc4005be48a7355c8cd010",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  finishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpoch),
				PreStoredInstances:      preInstances(testingutils.ProposerRunner(ks)),
				Duty:                    testingutils.TestingProposerDutyNextEpoch,
				PostDutyRunnerStateRoot: "5e25bba69a32033927e36cceb15b761f4d7c89bf991efbdda657d4a5da77a90d",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  finishRunner(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				PreStoredInstances:      preInstances(testingutils.AttesterRunner(ks)),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "fcd05356f14f26ac9152ff37de7f4eb2e7d4d370e4c68b319acff2b576668680",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
		},
	}
}
