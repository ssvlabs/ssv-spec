package postconsensus

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	ssz "github.com/ferranbt/fastssz"
)

func getSSZRootNoError(obj ssz.HashRoot) spec.Root {
	r, _ := obj.HashTreeRoot()
	return r
}

func decideRunner(r ssv.Runner, duty *types.Duty, decidedValue *types.ConsensusData) ssv.Runner {
	r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
	r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
		r.GetBaseRunner().QBFTController.GetConfig(),
		r.GetBaseRunner().Share,
		r.GetBaseRunner().QBFTController.Identifier,
		qbft.FirstHeight)
	r.GetBaseRunner().State.RunningInstance.State.Decided = true
	r.GetBaseRunner().State.DecidedValue = decidedValue
	r.GetBaseRunner().QBFTController.StoredInstances[0] = r.GetBaseRunner().State.RunningInstance
	r.GetBaseRunner().QBFTController.Height = qbft.FirstHeight
	return r
}

// ValidMessage tests a valid SignedPartialSignatureMessage with multi PartialSignatureMessages
func ValidMessage() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus valid msg",
		Tests: []*tests.MsgProcessingSpecTest{
			//{
			//	Name:   "sync committee aggregator selection proof",
			//	Runner: testingutils.SyncCommitteeContributionRunner(ks),
			//	Duty:   testingutils.TestingSyncCommitteeContributionDuty,
			//	Messages: []*types.SSVMessage{
			//		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			//	},
			//	PostDutyRunnerStateRoot: "211383dd16eb1fcc3264ec5d41d930cb899f7762da5b3067611549a7e6c8cb76",
			//	OutputMessages: []*ssv.SignedPartialSignatureMessage{
			//		testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
			//	},
			//},
			//{
			//	Name:   "aggregator selection proof",
			//	Runner: testingutils.AggregatorRunner(ks),
			//	Duty:   testingutils.TestingAggregatorDuty,
			//	Messages: []*types.SSVMessage{
			//		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			//	},
			//	PostDutyRunnerStateRoot: "9fde6ba2724a004c48f9c2b633797f14561370a62359cd555d084892ee1bf5c3",
			//	OutputMessages: []*ssv.SignedPartialSignatureMessage{
			//		testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
			//	},
			//},
			//{
			//	Name:   "randao",
			//	Runner: testingutils.ProposerRunner(ks),
			//	Duty:   testingutils.TestingProposerDuty,
			//	Messages: []*types.SSVMessage{
			//		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			//	},
			//	PostDutyRunnerStateRoot: "7ef3001d70abe0e85ad00653810e7395ffe9dfae09b52df45e5fa1770f252173",
			//	OutputMessages: []*ssv.SignedPartialSignatureMessage{
			//		testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
			//	},
			//},
			{
				Name: "attester",
				Runner: decideRunner(
					testingutils.AttesterRunner(ks),
					testingutils.TestingAttesterDuty,
					testingutils.TestAttesterConsensusData,
				),
				Duty: testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
				},
				PostDutyRunnerStateRoot: "d18cac07942b00be832092a6d548aaf0468d581450bc5eb443ff213f10168957",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []spec.Root{},
				DontStartDuty:           true,
			},
			//{
			//	Name:   "sync committee",
			//	Runner: testingutils.SyncCommitteeRunner(ks),
			//	Duty:   testingutils.TestingSyncCommitteeDuty,
			//	Messages: []*types.SSVMessage{
			//		testingutils.SSVMsgSyncCommittee(nil, testingutils.PreConsensusFailedMsg(ks.Shares[1], 1)),
			//	},
			//	PostDutyRunnerStateRoot: "9d06c3b83aee2bf5723ac0a19fdb9d011eeb87694575e27cc3775a8772eedbfa",
			//	OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
			//	ExpectedError:           "no pre consensus sigs required for sync committee role",
			//},
		},
	}
}
