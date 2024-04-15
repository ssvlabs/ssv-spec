package newduty

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FirstHeight tests a valid start duty at slot 0
func FirstHeight() tests.SpecTest {

	panic("implement me")

	ks := testingutils.Testing4SharesSet()
	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty first height",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:   "sync committee aggregator",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDutyFirstSlot,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusCustomSlotContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, 0),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDutyFirstSlot,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusCustomSlotSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, 0),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDutyFirstSlot,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:           "attester and sync committee",
				Runner:         testingutils.ClusterRunner(ks),
				Duty:           &testingutils.TestingAttesterDutyFirstSlot,
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}
}
