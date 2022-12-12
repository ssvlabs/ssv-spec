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

	// TODO: check error
	// nolint
	finishRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.FirstHeight
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
				PostDutyRunnerStateRoot: "c00253596603e226341c1d54b10f5a23a8d3055a83efa48a1cdc609ed5fe828a",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  finishRunner(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "02274203580a12cc3051ab79b90fc1514f788653791269cf681158f344689a56",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  finishRunner(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDutyNextEpoch),
				Duty:                    testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "1021eb2f5a418074ba6ef7d041ca359949c524665ce9579925b8752551c00632",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  finishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpoch),
				Duty:                    testingutils.TestingProposerDutyNextEpoch,
				PostDutyRunnerStateRoot: "970b351ea007de26b72ebdbc108431fa8429f8f2d7a30e6258beabe754dd31d8",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  finishRunner(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "9302df952cc59718f2b7d3edc59be07a4564b2a34846f84000bd56e076f267b1",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			},
		},
	}
}
