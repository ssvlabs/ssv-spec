package newduty

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
				PostDutyRunnerStateRoot: "cb45390b7f7d56a825613fa9e6c3fbc98fc565d8bd9c25d31ef01937dfff72b8",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  testingutils.SyncCommitteeRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "4ecc5fa28895cda24f3c9489c891c69431ad1c9aedf66d27f87528494b313e19",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "d898d4dff1fc0e6fd56ff075adc9766247cc915530e4fabf70324922ad7454ec",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDuty,
				PostDutyRunnerStateRoot: "d6be557a9c94566c60246966c15eca49d458117dce1ba381660d80360a607354",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  testingutils.AttesterRunner(ks),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "d00a1795c598a53727f7b683d54b2299289dd879c82198123556537e781fcede",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
		},
	}
}
