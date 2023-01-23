package newduty

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Valid tests a valid start duty
func Valid() *MultiStartNewRunnerDutySpecTest {
	ks := testingutils.Testing4SharesSet()
	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty valid",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
				PostDutyRunnerStateRoot: "93f63d09f9a12244c8ca751934450e7be81f89eb876dc76119ae7a7153ed405c",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  testingutils.SyncCommitteeRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "2c4a3a6aa312f25fdb6eef98fc9036eb3db007c6a4a05f3e6137737189acce49",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "57a6fb4d6bf799353c3c0980e3eb31ff32f0c100ebe04f6025dc06fc983ef52b",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDuty,
				PostDutyRunnerStateRoot: "65cd178a17356e48a1d7de571c0d56f791f6c4654b5b03e385adea510f7c55d1",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  testingutils.AttesterRunner(ks),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "27a84d7d33b8851a5b0e241280ddb160049768179c74295cb9c35fb875fe44e7",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
		},
	}
}
