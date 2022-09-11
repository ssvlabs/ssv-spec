package newduty

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Finished tests a valid start duty after finished prev
func Finished() *MultiStartNewRunnerDutySpecTest {
	ks := testingutils.Testing4SharesSet()

	finishRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.StartNewDuty(duty)
		r.GetBaseRunner().State.RunningInstance = &qbft.Instance{State: &qbft.State{Decided: true}}
		r.GetQBFTController().StoredInstances[0] = &qbft.Instance{State: &qbft.State{Decided: true}}
		r.GetQBFTController().Height = qbft.FirstHeight
		r.GetBaseRunner().State.Finished = true
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty finished",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  finishRunner(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionNexEpochDuty),
				Duty:                    testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "cc33ce8878818571f286ebd4d42d72b63608b5ef342aa7108ba6838512bf440f",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  finishRunner(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "f40c3a7ec6f4645793aaf5b17c0b85afbe7bd09570eda27de32d3218645d1489",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  finishRunner(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDutyNextEpoch),
				Duty:                    testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "42ad2eaefc8ff9eed997300c3f52bc50b57429cfe5e3551f45f2a959c6bdec7f",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  finishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpoch),
				Duty:                    testingutils.TestingProposerDutyNextEpoch,
				PostDutyRunnerStateRoot: "8fccf013fdb9b1fd99f0b9d0c31dca60d662c332bda9719ffddfadac455fcfca",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
					testingutils.PreConsensusRandaoNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  finishRunner(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "b6c4f69463e6100ea5480e3711654a34fa9004e0703442ee8743d090ce5ec0df",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
		},
	}
}
