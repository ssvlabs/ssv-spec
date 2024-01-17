package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidThenQuorum  tests a runner receiving an invalid message then a valid quorum, terminating successfully
func InvalidThenQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	expectedError := "failed processing post consensus message: invalid post-consensus message: could not verify Beacon partial Signature: wrong signature"
	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus invalid then quorum",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee contribution",
				Runner: decideRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty,
					testingutils.TestSyncCommitteeContributionConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusWrongSigSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),

					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks)),
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks)),
				},
				PostDutyRunnerStateRoot: invalidThenQuorumSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     invalidThenQuorumSyncCommitteeContributionSC().ExpectedState,
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
				DontStartDuty: true,
				ExpectedError: expectedError,
			},
			{
				Name: "sync committee",
				Runner: decideRunner(
					testingutils.SyncCommitteeRunner(ks),
					&testingutils.TestingSyncCommitteeDuty,
					testingutils.TestSyncCommitteeConsensusData,
				),
				Duty: &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusWrongSigSyncCommitteeMsg(ks.Shares[1], 1)),

					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3)),
				},
				PostDutyRunnerStateRoot: invalidThenQuorumSyncCommitteeSC().Root(),
				PostDutyRunnerState:     invalidThenQuorumSyncCommitteeSC().ExpectedState,
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks)),
				},
				DontStartDuty: true,
				ExpectedError: expectedError,
			},
			{
				Name: "aggregator",
				Runner: decideRunner(
					testingutils.AggregatorRunner(ks),
					&testingutils.TestingAggregatorDuty,
					testingutils.TestAggregatorConsensusData,
				),
				Duty: &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusWrongSigAggregatorMsg(ks.Shares[1], 1)),

					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
				},
				PostDutyRunnerStateRoot: invalidThenQuorumAggregatorSC().Root(),
				PostDutyRunnerState:     invalidThenQuorumAggregatorSC().ExpectedState,
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks)),
				},
				DontStartDuty: true,
				ExpectedError: expectedError,
			},
			{
				Name: "attester",
				Runner: decideRunner(
					testingutils.AttesterRunner(ks),
					&testingutils.TestingAttesterDuty,
					testingutils.TestAttesterConsensusData,
				),
				Duty: &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusWrongSigAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),

					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)),
					testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)),
				},
				PostDutyRunnerStateRoot: invalidThenQuorumAttesterSC().Root(),
				PostDutyRunnerState:     invalidThenQuorumAttesterSC().ExpectedState,
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAttestation(ks)),
				},
				DontStartDuty: true,
				ExpectedError: expectedError,
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
			Messages: []*types.SSVMessage{
				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongSigProposerMsgV(ks.Shares[1], 1, version)),

				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version)),
			},
			PostDutyRunnerStateRoot: invalidThenQuorumProposerSC(version).Root(),
			PostDutyRunnerState:     invalidThenQuorumProposerSC(version).ExpectedState,
			OutputMessages:          []*types.SignedPartialSignatureMessage{},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
			},
			DontStartDuty: true,
			ExpectedError: expectedError,
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
			Messages: []*types.SSVMessage{
				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongSigProposerMsgV(ks.Shares[1], 1, version)),

				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version)),
			},
			PostDutyRunnerStateRoot: invalidThenQuorumBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     invalidThenQuorumBlindedProposerSC(version).ExpectedState,
			OutputMessages:          []*types.SignedPartialSignatureMessage{},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
			},
			DontStartDuty: true,
			ExpectedError: expectedError,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
