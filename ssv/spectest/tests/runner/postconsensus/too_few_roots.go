package postconsensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// TooFewRoots tests a valid SignedPartialSignatureMessage with too few roots
func TooFewRoots() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	err := "failed processing post consensus message: invalid post-consensus message: SignedPartialSignatureMessage invalid: no PartialSignatureMessages messages"
	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus too few roots",
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
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionTooFewRootsMsg(ks.Shares[1], 1, ks)),
				},
				PostDutyRunnerStateRoot: "fb3ace1aa8af8989d9cd0829d375df27f413b2e28ac77ef7fe75df7de4718d51",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing post consensus message: invalid post-consensus message: wrong expected roots count",
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
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeTooFewRootsMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "a804b056eec86ed6af5f76ec2fe71d5c448ccef4eeea41e04018d000eece2032",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           err,
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
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerTooFewRootsMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "80841cc4f4befad1569c1b13c95f2512e45d5c6408d6c8a0a13114a46f30de10",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name: "proposer (blinded block)",
				Runner: decideRunner(
					testingutils.ProposerBlindedBlockRunner(ks),
					testingutils.TestingProposerDuty,
					testingutils.TestProposerBlindedBlockConsensusData,
				),
				Duty: testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerTooFewRootsMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "a2e539418ce3768e21786ab7a7377a100d6064d6a02fb6c0b83aad2f7e12485c",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           err,
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
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorTooFewRootsMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "e23696bceac6737e9e9bdd7be9d33f61c3f03308d91eea4225a9f8614132e5e1",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           err,
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
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationTooFewRootsMsg(ks.Shares[1], 1, qbft.FirstHeight)),
				},
				PostDutyRunnerStateRoot: "78ccc8a259311030f8c162aa66bf86dcdae0ebe3abf99ffdb61df3050af66950",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
		},
	}
}
