package preconsensus

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnorderedExpectedRoots tests expected roots to match but out of order, should return error
func UnorderedExpectedRoots() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus unordered expected roots",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusWrongOrderContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "cb45390b7f7d56a825613fa9e6c3fbc98fc565d8bd9c25d31ef01937dfff72b8",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing sync committee selection proof message: invalid pre-consensus message: wrong pre consensus signing root",
			},
			{
				Name:                    "aggregator selection proof",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    testingutils.TestingAggregatorDuty,
				Messages:                []*types.SSVMessage{},
				PostDutyRunnerStateRoot: "d898d4dff1fc0e6fd56ff075adc9766247cc915530e4fabf70324922ad7454ec",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "randao",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDuty,
				Messages:                []*types.SSVMessage{},
				PostDutyRunnerStateRoot: "d6be557a9c94566c60246966c15eca49d458117dce1ba381660d80360a607354",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
		},
	}
}
