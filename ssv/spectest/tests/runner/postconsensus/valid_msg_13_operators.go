package postconsensus

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// ValidMessage13Operators tests a valid SignedPartialSignatureMessage with multi PartialSignatureMessages with 13 operators
func ValidMessage13Operators() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing13SharesSet()

	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus valid msg 13 operators",
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
				PostDutyRunnerStateRoot: "32d5b5348b731095647a14d5bf0dbb35d701146efda281ee36753a4ca767b159",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
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
				PostDutyRunnerStateRoot: "0bb1fd7874815822d0250ffc3651b0e4198c6fbfebb64cf6453b0e1f65c3b68c",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
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
				PostDutyRunnerStateRoot: "06d51ff91682182e2e2d134aa3dabc40bcec9d30f9d8867f6c90f6624b7f4b45",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
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
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "e302be9084870a35e843540c99fac544a3556f2b95cc9a773a453691eceb2633",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
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
				PostDutyRunnerStateRoot: "9085e6dfb3a09bbfcdc11a5c6c48a3309720a669d3dbaebab4730cf2c7f8acd6",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
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
				PostDutyRunnerStateRoot: "168b5ef4821d9022b05bf3d0c0c82a8bb462778ec0bf07794ee876a7a25e507e",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
			},
		},
	}
}
