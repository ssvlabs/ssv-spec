package newduty

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotDecided tests starting duty before finished or decided
func NotDecided() *MultiStartNewRunnerDutySpecTest {
	ks := testingutils.Testing4SharesSet()

	preInstances := func(r ssv.Runner) []*qbft.Instance {
		return []*qbft.Instance{qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)}
	}

	// TODO: check error
	// nolint
	startRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
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
		Name: "new duty not decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  startRunner(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionNexEpochDuty),
				PreStoredInstances:      preInstances(testingutils.SyncCommitteeContributionRunner(ks)),
				Duty:                    testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "c8c2a180df39aec90cec21a5349ecef71fbe634787cf35d74566580ca3351baf",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "consensus on duty is running",
			},
			{
				Name:                    "sync committee",
				Runner:                  startRunner(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				PreStoredInstances:      preInstances(testingutils.SyncCommitteeRunner(ks)),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "9b294facd2b7a886f1c7c9e08ebc0381b29d70cd32993c682473d5a34265e189",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "consensus on duty is running",
			},
			{
				Name:                    "aggregator",
				Runner:                  startRunner(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDutyNextEpoch),
				PreStoredInstances:      preInstances(testingutils.AggregatorRunner(ks)),
				Duty:                    testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "b0363a34bd5936ae8bd5c9469cd3671fcd2e3bf46fdb8b0a357d0ca7d7c76174",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "consensus on duty is running",
			},
			{
				Name:                    "proposer",
				Runner:                  startRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpoch),
				PreStoredInstances:      preInstances(testingutils.ProposerRunner(ks)),
				Duty:                    testingutils.TestingProposerDutyNextEpoch,
				PostDutyRunnerStateRoot: "c72233eb1d9ce1431ee821ca72f2cbe228e2a6f83a16db0507bb4b9458f9ccd3",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "consensus on duty is running",
			},
			{
				Name:                    "attester",
				Runner:                  startRunner(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				PreStoredInstances:      preInstances(testingutils.AttesterRunner(ks)),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "6f4f489101a2270e5defa1fba2467813a109e028eed95cdd9e4d079d360b0b0f",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "consensus on duty is running",
			},
		},
	}
}
