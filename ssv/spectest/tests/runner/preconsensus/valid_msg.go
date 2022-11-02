package preconsensus

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidMessage tests a valid SignedPartialSignature with multi PartialSignatures
func ValidMessage() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus valid msg",
		Tests: []*tests.MsgProcessingSpecTest{
			//{
			//	Name:   "sync committee aggregator selection proof",
			//	Runner: testingutils.SyncCommitteeContributionRunner(ks),
			//	Duty:   testingutils.TestingSyncCommitteeContributionDuty,
			//	Messages: []*types.Message{
			//		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), types.PartialContributionProofSignatureMsgType),
			//	},
			//	PostDutyRunnerStateRoot: "3ab9ba2158fe8b2706b63ec95edf26eb556ad44ec509e1040000ab37fd89da72",
			//	OutputMessages: []*ssv.SignedPartialSignature{
			//		testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
			//	},
			//},
			//{
			//	Name:   "aggregator selection proof",
			//	Runner: testingutils.AggregatorRunner(ks),
			//	Duty:   testingutils.TestingAggregatorDuty,
			//	Messages: []*types.Message{
			//		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), types.PartialSelectionProofSignatureMsgType),
			//	},
			//	PostDutyRunnerStateRoot: "212c6c21628e2d13a407183badbee87e06c423b981c37a22263ce94d7100f370",
			//	OutputMessages: []*ssv.SignedPartialSignature{
			//		testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
			//	},
			//},
			//{
			//	Name:   "randao",
			//	Runner: testingutils.ProposerRunner(ks),
			//	Duty:   testingutils.TestingProposerDuty,
			//	Messages: []*types.Message{
			//		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1), types.PartialRandaoSignatureMsgType),
			//	},
			//	PostDutyRunnerStateRoot: "955508489a4ca9abad520e72aa4a1f27c2f8631979b2f5b3a926572543995bb5",
			//	OutputMessages: []*ssv.SignedPartialSignature{
			//		testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
			//	},
			//},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAttester(nil, testingutils.PreConsensusFailedMsg(ks.Shares[1], 1), types.PartialContributionProofSignatureMsgType),
				},
				PostDutyRunnerStateRoot: "d2ff2fad2e2b99af0f3e1f0378384ea9e89815d4dc4ca93f8754583a4a019b07",
				OutputMessages:          []*ssv.SignedPartialSignature{},
				ExpectedError:           "no pre consensus sigs required for attester role",
			},
			//{
			//	Name:   "sync committee",
			//	Runner: testingutils.SyncCommitteeRunner(ks),
			//	Duty:   testingutils.TestingSyncCommitteeDuty,
			//	Messages: []*types.Message{
			//		testingutils.SSVMsgSyncCommittee(nil, testingutils.PreConsensusFailedMsg(ks.Shares[1], 1), types.PartialContributionProofSignatureMsgType),
			//	},
			//	PostDutyRunnerStateRoot: "971bbb10725791ebfd0a06f9817797feac2b81b17babb538b3c3f34421341cda",
			//	OutputMessages:          []*ssv.SignedPartialSignature{},
			//	ExpectedError:           "no pre consensus sigs required for sync committee role",
			//},
		},
	}
}
