package postconsensus

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Quorum7Operators  tests a quorum of valid SignedPartialSignatureMessage 7 operators
func Quorum7Operators() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing7SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus quorum 7 operators",
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
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[4], 4, ks)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[5], 5, ks)),
				},
				PostDutyRunnerStateRoot: "34e5ae4a494d54cc2042a1c44f388e54d896f19cf85f76474ed3caab55742aa5",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
				DontStartDuty: true,
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
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[5], 5)),
				},
				PostDutyRunnerStateRoot: "3f50c8d9f05fe6a3dca83cb1f2aa851226d3a8ad43c0e6894438ee9ed7748944",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks)),
				},
				DontStartDuty: true,
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
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[5], 5)),
				},
				PostDutyRunnerStateRoot: "e7f81015870281deade1f08c9b4fc0a906396c58b4e6933a51cfd5615149de3f",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
				},
				DontStartDuty: true,
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
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[5], 5)),
				},
				PostDutyRunnerStateRoot: "c75fc4f697839c65065b76450cf6aa0a44fce43cdb457a1bebc5712fcb4a15a7",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
				},
				DontStartDuty: true,
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
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[5], 5)),
				},
				PostDutyRunnerStateRoot: "3c7c4a1cb25043b9032e3420391d3c5288ecd4bc62e1e2148c86c2c948c384b3",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks)),
				},
				DontStartDuty: true,
			},
			{
				Name: "attester",
				Runner: decideRunner(
					testingutils.AttesterRunner(ks),
					testingutils.TestingAttesterDuty,
					testingutils.TestAttesterConsensusData,
				),
				Duty: testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[4], 4, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[5], 5, qbft.FirstHeight)),
				},
				PostDutyRunnerStateRoot: "7ba346220ec36ddeea29938e7b92112fc394350e8aedf165fc1cb9db8dc696e0",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedAttestation(ks)),
				},
				DontStartDuty: true,
			},
		},
	}
}
