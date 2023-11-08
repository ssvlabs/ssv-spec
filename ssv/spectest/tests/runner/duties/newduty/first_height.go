package newduty

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FirstHeight tests a valid start duty at slot 0
func FirstHeight() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty first height",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:   "sync committee aggregator",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDutyFirstSlot,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusCustomSlotContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, 0),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:           "sync committee",
				Runner:         testingutils.SyncCommitteeRunner(ks),
				Duty:           &testingutils.TestingSyncCommitteeDutyFirstSlot,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDutyFirstSlot,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusCustomSlotSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, 0),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDutyFirstSlot,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:           "attester",
				Runner:         testingutils.AttesterRunner(ks),
				Duty:           &testingutils.TestingAttesterDutyFirstSlot,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
		},
	}
}
