package newduty

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotDecided tests starting duty before finished or decided
func NotDecided() *MultiStartNewRunnerDutySpecTest {
	ks := testingutils.Testing4SharesSet()

	startRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.StartNewDuty(duty)
		r.GetBaseRunner().State.RunningInstance = &qbft.Instance{State: &qbft.State{Decided: false}}
		r.GetQBFTController().StoredInstances[0] = &qbft.Instance{State: &qbft.State{Decided: false}}
		r.GetQBFTController().Height = qbft.FirstHeight
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty not decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  startRunner(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionNexEpochDuty),
				Duty:                    testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "e7882af6b88191ff7661b9b824b916fc9210312765c901049ae0bc5d584b5082",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "consensus on duty is running",
			},
			{
				Name:                    "sync committee",
				Runner:                  startRunner(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "4fd4745540ea47ec69577474065129b6b2717f0122b38e7a25828eeb00ee23d6",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "consensus on duty is running",
			},
			{
				Name:                    "aggregator",
				Runner:                  startRunner(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDutyNextEpoch),
				Duty:                    testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "a50a96edc2d373b124ff8749794b680ead49d1e90d1de6482ac7796494206d23",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "consensus on duty is running",
			},
			{
				Name:                    "proposer",
				Runner:                  startRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpoch),
				Duty:                    testingutils.TestingProposerDutyNextEpoch,
				PostDutyRunnerStateRoot: "e80c00a1afbfe1fc43ab101ac2483873fbdcbc4555cfbea59f1c51484ce29948",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "consensus on duty is running",
			},
			{
				Name:                    "attester",
				Runner:                  startRunner(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				Duty:                    testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "3c86e7cec5c28aed9202dab84d356e6655bdafc9de13ee89a51a444b1722d210",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "consensus on duty is running",
			},
		},
	}
}
