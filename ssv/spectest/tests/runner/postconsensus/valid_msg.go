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
			{
				Name: "sync committee contribution",
				Runner: decideRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					testingutils.TestingSyncCommitteeContributionDuty,
					testingutils.TestSyncCommitteeContributionConsensusData,
				),
				Duty: testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
				},
				PostDutyRunnerStateRoot: "fc2ce4bab95761ac0416c0fe8443defa829d9e033c905e37a539138abb700628",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []spec.Root{},
				DontStartDuty:           true,
			},
			{
				Name: "sync committee",
				Runner: decideRunner(
					testingutils.SyncCommitteeRunner(ks),
					testingutils.TestingSyncCommitteeDuty,
					testingutils.TestSyncCommitteeConsensusData,
				),
				Duty: testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "b10fe3e270c3ed59a89d1b1a5017a93b3c44665957898a1556cc1679e0921687",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []spec.Root{},
				DontStartDuty:           true,
			},
			{
				Name: "proposer",
				Runner: decideRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDuty,
					testingutils.TestProposerConsensusData,
				),
				Duty: testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "90d9f1b97b673119b32c7f0665f4ff8e8e02fd3ae34fdc902bef563e9e96519d",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []spec.Root{},
				DontStartDuty:           true,
			},
			{
				Name: "aggregator",
				Runner: decideRunner(
					testingutils.AggregatorRunner(ks),
					testingutils.TestingAggregatorDuty,
					testingutils.TestAggregatorConsensusData,
				),
				Duty: testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "f331cb67db71d735e60d8e0c5b23b2e54bef5054679ac84634a66d62bc0edca1",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []spec.Root{},
				DontStartDuty:           true,
			},
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
		},
	}
}
