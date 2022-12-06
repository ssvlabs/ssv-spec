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
				PostDutyRunnerStateRoot: "660668ba29134cf00a00c73221533f9771e61a5fd0df2545e37406695860b532",
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
				PostDutyRunnerStateRoot: "12620c4d79ec66cf9b4058f707f210373fab6e5dad603696c7a8294dc82385a3",
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
				PostDutyRunnerStateRoot: "ffcdadcfa1ebc4b750102c854a41b3355f32316187992666bf0c568b026df9c0",
				DontStartDuty:           true,
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "failed processing randao message: invalid pre-consensus message: no running duty",
			},
		},
	}
}
