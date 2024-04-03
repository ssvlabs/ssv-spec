package postconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Quorum13Operators  tests a quorum of valid SignedPartialSignatureMessage 13 operators
func Quorum13Operators() tests.SpecTest {
	ks := testingutils.Testing13SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus quorum 13 operators",
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
					testingutils.SSVMsgSyncCommitteeContribution(8, ks.NetworkKeys[8], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[8], 8)),
					testingutils.SSVMsgSyncCommitteeContribution(9, ks.NetworkKeys[9], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[9], 9)),
				},
				PostDutyRunnerStateRoot: "2813af9cbe896e94ab82dc427b8ae27cdc7aae7d5ea55ad3ab1e605839ad4c27",
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
					testingutils.SSVMsgSyncCommittee(8, ks.NetworkKeys[8], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[8], 8)),
					testingutils.SSVMsgSyncCommittee(9, ks.NetworkKeys[9], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[9], 9)),
				},
				PostDutyRunnerStateRoot: "df4f320a4721aaa6e67c85aa7048442db6d51d1271cde130194ef657faa893b2",
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
					testingutils.SSVMsgProposer(8, ks.NetworkKeys[8], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[8], 8, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(9, ks.NetworkKeys[9], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[9], 9, spec.DataVersionDeneb)),
				},
				PostDutyRunnerStateRoot: "68bf9032da3fe24ddf90b3bba686a684cb7d7ba0a13b86e54c5b096916be8d02",
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
					testingutils.SSVMsgProposer(8, ks.NetworkKeys[8], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[8], 8, spec.DataVersionDeneb)),
					testingutils.SSVMsgProposer(9, ks.NetworkKeys[9], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[9], 9, spec.DataVersionDeneb)),
				},
				PostDutyRunnerStateRoot: "404287c64ae601ac20744783becd0081ff25da8b58592daac39365a8bbdb5e4a",
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
					testingutils.SSVMsgAggregator(8, ks.NetworkKeys[8], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[8], 8)),
					testingutils.SSVMsgAggregator(9, ks.NetworkKeys[9], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[9], 9)),
				},
				PostDutyRunnerStateRoot: "9c2b6bef7eccfc15895d9db88b5dcb01a17b6300b6efbdbf0d0ba263d7655899",
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
					testingutils.SSVMsgAttester(8, ks.NetworkKeys[8], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[8], 8, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(9, ks.NetworkKeys[9], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[9], 9, qbft.FirstHeight)),
				},
				PostDutyRunnerStateRoot: "f1dfaa89b18ddf29bdbf67698b986dd29c43c7b8e87981de5540394cab3e087a",
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAttestation(ks)),
				},
				DontStartDuty: true,
			},
		},
	}
}
