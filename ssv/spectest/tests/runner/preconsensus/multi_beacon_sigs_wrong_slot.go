package preconsensus

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MultiBeaconSigsWrongSlot tests SignedPartialSignatures with multi PartialSignatures where one has slot != msg.Slot
func MultiBeaconSigsWrongSlot() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus multi msg wrong slot",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusWrongMsgSlotContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), types.PartialContributionProofSignatureMsgType),
				},
				PostDutyRunnerStateRoot: "f7d21e5fafab57daf6fe0a0fb9efc50fea94e40c87394aa58ab5b9c3569e4042",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing sync committee selection proof message: invalid pre-consensus message: wrong slot",
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusWrongMsgSlotSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), types.PartialSelectionProofSignatureMsgType),
				},
				PostDutyRunnerStateRoot: "da498b2ca86a535b3b879ddcef5d9de924162c478694129a6a21857641e6031a",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing selection proof message: invalid pre-consensus message: wrong slot",
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoWrongSlotMsg(ks.Shares[1], 1), types.PartialRandaoSignatureMsgType),
				},
				PostDutyRunnerStateRoot: "583f8c730925b043befceb4e871a9cf6c28b3f0e5c1b173ca0d869e958e8445f",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "failed processing randao message: invalid pre-consensus message: wrong slot",
			},
		},
	}
}
