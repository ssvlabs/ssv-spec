package preconsensus

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongExpectedRootsCount tests len(expectedRootObjs) != len(signedMsg.Message.Messages)
func WrongExpectedRootsCount() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus wrong expected roots count",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusWrongCountContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "a00545e4062f9a566febdca8222f89d2cea35bcfb15de4fb86f15b8746e63e56",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing sync committee selection proof message: invalid pre-consensus message: wrong expected roots count",
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusMultiSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "932b78133774f19a0f9fe303243ac3ca69946c925280c1fdf25c22e11fa74e7c",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing selection proof message: invalid pre-consensus message: wrong expected roots count",
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMultiMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "ca3d758a37f4448b654c844b2990ea8fe705920ee31b9732ed7bcb92ac7e5400",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing randao message: invalid pre-consensus message: wrong expected roots count",
			},
		},
	}
}
