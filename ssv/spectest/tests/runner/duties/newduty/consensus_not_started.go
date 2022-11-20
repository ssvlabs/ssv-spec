package newduty

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ConsensusNotStarted tests starting duty after prev already started but for some duties' consensus didn't start because pre-consensus didnt get quorum (different duties will enable starting a new duty)
func ConsensusNotStarted() *MultiStartNewRunnerDutySpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	startRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty consensus not started",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  startRunner(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
				PostDutyRunnerStateRoot: "6229a8a7cabfac839fcd61eb0a00e67d5b150f51d0394c45d5d6f0b55659d8e6",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				PreStoredInstances: []*qbft.Instance{},
			},
			{
				Name:                    "sync committee",
				Runner:                  startRunner(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "9b294facd2b7a886f1c7c9e08ebc0381b29d70cd32993c682473d5a34265e189",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				PreStoredInstances:      []*qbft.Instance{},
			},
			{
				Name:                    "aggregator",
				Runner:                  startRunner(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDuty),
				Duty:                    testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "319d3a46e9cd8cc942d9bdf531cb92207696c0ba3bc86199ddb30adb2c673d4a",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				PreStoredInstances: []*qbft.Instance{},
			},
			{
				Name:                    "proposer",
				Runner:                  startRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDuty),
				Duty:                    testingutils.TestingProposerDuty,
				PostDutyRunnerStateRoot: "1c3fc4711119e95a6980cd74e2b7c7c9ac3a1c4698ce1ded1a58bd7f4c4c6a44",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				PreStoredInstances: []*qbft.Instance{},
			},
			{
				Name:                    "attester",
				Runner:                  startRunner(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "6f4f489101a2270e5defa1fba2467813a109e028eed95cdd9e4d079d360b0b0f",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				PreStoredInstances:      []*qbft.Instance{},
			},
		},
	}
}
