package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostQuorum tests a msg received post post-consensus quorum
func PostQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedErr := "failed processing post consensus message: invalid post-consensus message: no running duty"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus post quorum",
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
				},
				PostDutyRunnerStateRoot: postQuorumSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     postQuorumSyncCommitteeContributionSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
				DontStartDuty: true,
				ExpectedError: expectedErr,
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
				},
				PostDutyRunnerStateRoot: postQuorumSyncCommitteeSC().Root(),
				PostDutyRunnerState:     postQuorumSyncCommitteeSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks)),
				},
				DontStartDuty: true,
				ExpectedError: expectedErr,
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
				},
				PostDutyRunnerStateRoot: postQuorumAggregatorSC().Root(),
				PostDutyRunnerState:     postQuorumAggregatorSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks)),
				},
				DontStartDuty: true,
				ExpectedError: expectedErr,
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
				},
				PostDutyRunnerStateRoot: postQuorumAttesterSC().Root(),
				PostDutyRunnerState:     postQuorumAttesterSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAttestation(ks)),
				},
				DontStartDuty: true,
				ExpectedError: expectedErr,
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name: fmt.Sprintf("proposer (%s)", version.String()),
			Runner: decideRunner(
				testingutils.ProposerRunner(ks),
				testingutils.TestingProposerDutyV(version),
				testingutils.TestProposerConsensusDataV(version),
			),
			Duty: testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version)),
				testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, version)),
			},
			PostDutyRunnerStateRoot: postQuorumProposerSC(version).Root(),
			PostDutyRunnerState:     postQuorumProposerSC(version).ExpectedState,
			OutputMessages:          []*types.PartialSignatureMessages{},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
			},
			DontStartDuty: true,
			ExpectedError: expectedErr,
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name: fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner: decideRunner(
				testingutils.ProposerBlindedBlockRunner(ks),
				testingutils.TestingProposerDutyV(version),
				testingutils.TestProposerBlindedBlockConsensusDataV(version),
			),
			Duty: testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version)),
				testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, version)),
			},
			PostDutyRunnerStateRoot: postQuorumBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     postQuorumBlindedProposerSC(version).ExpectedState,
			OutputMessages:          []*types.PartialSignatureMessages{},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBlindedBeaconBlockV(ks, version)),
			},
			DontStartDuty: true,
			ExpectedError: expectedErr,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
