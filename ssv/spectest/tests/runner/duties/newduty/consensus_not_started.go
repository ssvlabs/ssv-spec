package newduty

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
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
				PostDutyRunnerStateRoot: "93f63d09f9a12244c8ca751934450e7be81f89eb876dc76119ae7a7153ed405c",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  startRunner(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "2c4a3a6aa312f25fdb6eef98fc9036eb3db007c6a4a05f3e6137737189acce49",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  startRunner(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDuty),
				Duty:                    testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "57a6fb4d6bf799353c3c0980e3eb31ff32f0c100ebe04f6025dc06fc983ef52b",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  startRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDuty),
				Duty:                    testingutils.TestingProposerDuty,
				PostDutyRunnerStateRoot: "65cd178a17356e48a1d7de571c0d56f791f6c4654b5b03e385adea510f7c55d1",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  startRunner(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "27a84d7d33b8851a5b0e241280ddb160049768179c74295cb9c35fb875fe44e7",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
		},
	}
}
