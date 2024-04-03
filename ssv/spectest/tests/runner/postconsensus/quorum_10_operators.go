package postconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Quorum10Operators  tests a quorum of valid SignedPartialSignatureMessage 10 operators
func Quorum10Operators() tests.SpecTest {
	ks := testingutils.Testing10SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus quorum 10 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee contribution",
				Runner: decideRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty,
					testingutils.TestSyncCommitteeContributionConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgSyncCommitteeContribution(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgSyncCommitteeContribution(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgSyncCommitteeContribution(4, ks.NetworkKeys[4], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgSyncCommitteeContribution(5, ks.NetworkKeys[5], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[5], 5)),
					testingutils.SSVMsgSyncCommitteeContribution(6, ks.NetworkKeys[6], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[6], 6)),
					testingutils.SSVMsgSyncCommitteeContribution(7, ks.NetworkKeys[7], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[7], 7)),
				},
				PostDutyRunnerStateRoot: "6c04842d4562ae0796eebf1726bb1b2b626582c257aa04753a369dfb65512ec2",
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
				DontStartDuty: true,
			},
			{
				Name: "sync committee",
				Runner: decideRunner(
					testingutils.SyncCommitteeRunner(ks),
					&testingutils.TestingSyncCommitteeDuty,
					testingutils.TestSyncCommitteeConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgSyncCommittee(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgSyncCommittee(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgSyncCommittee(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgSyncCommittee(4, ks.NetworkKeys[4], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgSyncCommittee(5, ks.NetworkKeys[5], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[5], 5)),
					testingutils.SSVMsgSyncCommittee(6, ks.NetworkKeys[6], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[6], 6)),
					testingutils.SSVMsgSyncCommittee(7, ks.NetworkKeys[7], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[7], 7)),
				},
				PostDutyRunnerStateRoot: "4340646e0b4ad970e8154f961bc6a40600a9cb2f3250dc3b8c5d334697e65561",
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks)),
				},
				DontStartDuty: true,
			},
			{
				Name: "proposer",
				Runner: decideRunner(
					testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
					testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(5, ks.NetworkKeys[5], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[5], 5, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(6, ks.NetworkKeys[6], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[6], 6, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(7, ks.NetworkKeys[7], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[7], 7, spec.DataVersionDeneb)),
				},
				PostDutyRunnerStateRoot: "bd4046479e197ace695896fce712bd5c7b993cf3af717eba3030204bed729bd3",
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, spec.DataVersionDeneb)),
				},
				DontStartDuty: true,
			},
			{
				Name: "proposer (blinded block)",
				Runner: decideRunner(
					testingutils.ProposerBlindedBlockRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
					testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb),
				),
				Duty: testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(5, ks.NetworkKeys[5], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[5], 5, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(6, ks.NetworkKeys[6], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[6], 6, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(7, ks.NetworkKeys[7], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[7], 7, spec.DataVersionDeneb)),
				},
				PostDutyRunnerStateRoot: "d0498893b8a9386af88df849f3d88bcff88c495c4705ab74fb4560e16abde3fb",
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedBlindedBeaconBlockV(ks, spec.DataVersionDeneb)),
				},
				DontStartDuty: true,
			},
			{
				Name: "aggregator",
				Runner: decideRunner(
					testingutils.AggregatorRunner(ks),
					&testingutils.TestingAggregatorDuty,
					testingutils.TestAggregatorConsensusData,
				),
				Duty: &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgAggregator(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgAggregator(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgAggregator(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgAggregator(4, ks.NetworkKeys[4], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgAggregator(5, ks.NetworkKeys[5], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[5], 5)),
					testingutils.SSVMsgAggregator(6, ks.NetworkKeys[6], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[6], 6)),
					testingutils.SSVMsgAggregator(7, ks.NetworkKeys[7], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[7], 7)),
				},
				PostDutyRunnerStateRoot: "886abce3e475b83205828b08da281d830c2479e1b1547f9d7c82f9e1bae11b0a",
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks)),
				},
				DontStartDuty: true,
			},
			{
				Name: "attester",
				Runner: decideRunner(
					testingutils.AttesterRunner(ks),
					&testingutils.TestingAttesterDuty,
					testingutils.TestAttesterConsensusData,
				),
				Duty: &testingutils.TestingAttesterDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgAttester(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(4, ks.NetworkKeys[4], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[4], 4, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(5, ks.NetworkKeys[5], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[5], 5, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(6, ks.NetworkKeys[6], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[6], 6, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(7, ks.NetworkKeys[7], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[7], 7, qbft.FirstHeight)),
				},
				PostDutyRunnerStateRoot: "ff24bc80b8e109fd51e343e7da2c7b40802203e892ea52d96a25ebe7fceb9032",
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAttestation(ks)),
				},
				DontStartDuty: true,
			},
		},
	}
}
