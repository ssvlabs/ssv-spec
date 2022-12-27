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
				PostDutyRunnerStateRoot: "abe78d812d565ca7fe6feba32af83199977260ba39f1e91a409837f8247931fe",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  testingutils.SyncCommitteeRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "eefb398de44e5321980ca3474b4a480cf1e688e54ca9d5c74b5b2412f7693fc7",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "9a20d699ce5e1ac9106ce3b4684c9556c9959fdaa18264725278a944f8d708dd",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDuty,
				PostDutyRunnerStateRoot: "0c3e3ecd96f75003b56f4b8f02ef7da661674028a7a8c4c654e663971c53f94a",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  testingutils.AttesterRunner(ks),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "e2a107146e8fbbe70b617cc6ac31e1b43dda031ff3597a87ff8821dd49fe4bb7",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
		},
	}
}
