package newduty

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostWrongDecided tests starting a new duty after prev was decided wrongly (future decided)
func PostWrongDecided() *MultiStartNewRunnerDutySpecTest {
	ks := testingutils.Testing4SharesSet()

	preInstances := func(r ssv.Runner) []*qbft.Instance {
		current := qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		current.State.Decided = true

		higherDecided := qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			10)
		higherDecided.State.Decided = true
		higherDecided.State.DecidedValue = []byte{1, 2, 3, 4}
		return []*qbft.Instance{
			current,
			higherDecided,
		}
	}

	decideWrong := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningHeight = qbft.FirstHeight
		r.GetBaseRunner().QBFTController.Height = 10

		// support running tests without json
		for _, instance := range preInstances(r) {
			if err := r.GetBaseRunner().QBFTController.SaveInstance(instance); err != nil {
				panic(err)
			}
		}

		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post wrong decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  decideWrong(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
				PostDutyRunnerStateRoot: "eab49e351691a9b726c03747cfba57223c07499d52e439406bebf58cd33be345",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				PreStoredInstances: preInstances(testingutils.SyncCommitteeContributionRunner(ks)),
			},
			{
				Name:                    "sync committee",
				Runner:                  decideWrong(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "7c059a02f4bf14d99dc0a92986f6780accb429d32027ddb5c14866f140abeb2b",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				PreStoredInstances:      preInstances(testingutils.SyncCommitteeRunner(ks)),
			},
			{
				Name:                    "aggregator",
				Runner:                  decideWrong(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDuty),
				Duty:                    testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "1048f89f584cbc26bda20559613daa4df9f13b5dee4d43fbe549bfa65064add4",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				PreStoredInstances: preInstances(testingutils.AggregatorRunner(ks)),
			},
			{
				Name:                    "proposer",
				Runner:                  decideWrong(testingutils.ProposerRunner(ks), testingutils.TestingProposerDuty),
				Duty:                    testingutils.TestingProposerDuty,
				PostDutyRunnerStateRoot: "f3529ad09d7a06cef30214fc075eb3895d61fc59a6bd11081814d2a4749daf14",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				PreStoredInstances: preInstances(testingutils.ProposerRunner(ks)),
			},
			{
				Name:                    "attester",
				Runner:                  decideWrong(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "2f91e238dfdfe9be1f5cdc3f79e7423219efb7fb8ebd5bbca9610d4d58bca609",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				PreStoredInstances:      preInstances(testingutils.AttesterRunner(ks)),
			},
		},
	}
}
