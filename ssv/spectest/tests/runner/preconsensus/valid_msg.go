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
			//	PostDutyRunnerStateRoot: "18c1f02c9a7df3e71e10f03891ed8fd8b87bed43bead72f6b8ec0367d2bfbacf",
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
			//	PostDutyRunnerStateRoot: "18022e0768377879d579cc9da57e6fe315afac97a219f0b06c015601668f5a74",
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
			//	PostDutyRunnerStateRoot: "ec6371ba342e7747a78dd5dd2ead41221fe303a7f859a0b4f22f6bbb31063404",
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
				PostDutyRunnerStateRoot: "0a87eee70b6ee2583dd414d6f07f6f5c433975409896dafb51628f5e393a7458",
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
			//	PostDutyRunnerStateRoot: "ca7c1b5bb6a1b2b5d486da30bdef8a96db109cbe5691d1191a0671eaaafb5cf0",
			//	OutputMessages:          []*ssv.SignedPartialSignature{},
			//	ExpectedError:           "no pre consensus sigs required for sync committee role",
			//},
		},
	}
}
