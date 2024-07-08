package newduty

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FirstHeight tests a valid start duty at slot 0
func FirstHeight() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty first height",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:      "sync committee aggregator",
				Runner:    testingutils.SyncCommitteeContributionRunner(ks),
				Duty:      &testingutils.TestingSyncCommitteeContributionDutyFirstSlot,
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusCustomSlotContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, 0),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:      "aggregator",
				Runner:    testingutils.AggregatorRunner(ks),
				Duty:      &testingutils.TestingAggregatorDutyFirstSlot,
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusCustomSlotSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, 0),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:      "proposer",
				Runner:    testingutils.ProposerRunner(ks),
				Duty:      &testingutils.TestingProposerDutyFirstSlot,
				Threshold: ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:           "attester",
				Runner:         testingutils.CommitteeRunner(ks),
				Duty:           testingutils.TestingAttesterDutyFirstSlot,
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
			},
			{
				Name:           "sync committee",
				Runner:         testingutils.CommitteeRunner(ks),
				Duty:           testingutils.TestingSyncCommitteeDutyFirstSlot,
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
			},
			{
				Name:           "attester and sync committee",
				Runner:         testingutils.CommitteeRunner(ks),
				Duty:           testingutils.TestingSyncCommitteeDutyFirstSlot,
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}
}
