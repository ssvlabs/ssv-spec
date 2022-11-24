package preconsensus

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostFinish tests a msg received post runner finished
func PostFinish() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check errors
	// nolint
	finishRunner := func(runner ssv.Runner, duty *types.Duty) ssv.Runner {
		runner.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		runner.GetBaseRunner().State.Finished = true
		return runner
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus post finish",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee aggregator selection proof",
				Runner: finishRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					testingutils.TestingSyncCommitteeContributionDuty,
				),
				Duty: testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[4], ks.Shares[4], 4, 4)),
				},
				PostDutyRunnerStateRoot: "7c928d3cf4c6a4dde4333a247f8709e4a18578983ab68d2d0175b98fe2fc746d",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "failed processing sync committee selection proof message: invalid pre-consensus message: no running duty",
			},
			{
				Name: "aggregator selection proof",
				Runner: finishRunner(
					testingutils.AggregatorRunner(ks),
					testingutils.TestingAggregatorDuty,
				),
				Duty: testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[4], ks.Shares[4], 4, 4)),
				},
				PostDutyRunnerStateRoot: "c79b5b48b3efead65a9cc6106f0bcb1aba03940459ab4e1e64e0ba0a6eabaf47",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "failed processing selection proof message: invalid pre-consensus message: no running duty",
			},
			{
				Name: "randao",
				Runner: finishRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDuty,
				),
				Duty: testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[4], ks.Shares[4], 4, 4)),
				},
				PostDutyRunnerStateRoot: "559047e67bca877234e857dc79984ecf21026f31698badae16950af13ab14ccc",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "failed processing randao message: invalid pre-consensus message: no running duty",
			},
		},
	}
}
