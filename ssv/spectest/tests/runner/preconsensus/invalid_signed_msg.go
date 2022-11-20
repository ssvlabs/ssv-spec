package preconsensus

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidSignedMessage tests SignedPartialSignatureMessage.Validate() != nil
func InvalidSignedMessage() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	invaldiateMsg := func(msg *ssv.SignedPartialSignatureMessage) *ssv.SignedPartialSignatureMessage {
		msg.Signature = nil
		return msg
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus invalid signed msg",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, invaldiateMsg(testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
				},
				PostDutyRunnerStateRoot: "6229a8a7cabfac839fcd61eb0a00e67d5b150f51d0394c45d5d6f0b55659d8e6",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing sync committee selection proof message: invalid pre-consensus message: SignedPartialSignatureMessage invalid: SignedPartialSignatureMessage sig invalid",
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, invaldiateMsg(testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
				},
				PostDutyRunnerStateRoot: "319d3a46e9cd8cc942d9bdf531cb92207696c0ba3bc86199ddb30adb2c673d4a",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing selection proof message: invalid pre-consensus message: SignedPartialSignatureMessage invalid: SignedPartialSignatureMessage sig invalid",
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, invaldiateMsg(testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
				},
				PostDutyRunnerStateRoot: "1c3fc4711119e95a6980cd74e2b7c7c9ac3a1c4698ce1ded1a58bd7f4c4c6a44",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing randao message: invalid pre-consensus message: SignedPartialSignatureMessage invalid: SignedPartialSignatureMessage sig invalid",
			},
		},
	}
}
