package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicateMsgDifferentRootsThenQuorum tests that it can successfully reach a quorum and end the duty even after receiving the following error:
// - a duplicate SignedPartialSignatureMessage (from same signer) but with different roots
func DuplicateMsgDifferentRootsThenQuorum() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	expectedError := "failed processing post consensus message: invalid post-consensus message: wrong signing root"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus duplicate msg different root then quorum",
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusWrongSyncCommitteeContributionMsg(ks.Shares[1], 1, ks))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks))),
				},
				PostDutyRunnerStateRoot: duplicateMsgDifferentRootsThenQuorumSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     duplicateMsgDifferentRootsThenQuorumSyncCommitteeContributionSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
				DontStartDuty: true,
				ExpectedError: expectedError,
			},
		},
	}

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name: fmt.Sprintf("aggregator (%s)", version.String()),
			Runner: decideRunner(
				testingutils.AggregatorRunner(ks),
				testingutils.TestingAggregatorDuty(version),
				testingutils.TestAggregatorConsensusData(version),
			),
			Duty: testingutils.TestingAggregatorDuty(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusWrongAggregatorMsg(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3, version))),
			},
			OutputMessages: []*types.PartialSignatureMessages{},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks, version)),
			},
			DontStartDuty: true,
			ExpectedError: expectedError,
		},
		)
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{

			{
				Name: fmt.Sprintf("attester (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongAttestationMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingAttestationResponseBeaconObject(ks, version)),
				},
				DontStartDuty: true,
				// No error is expected for this duty because it don't overwrite (in state) the previous message
			},
			{
				Name: fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingSyncCommitteeDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingSyncCommitteeDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongSyncCommitteeMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks, version)),
				},
				DontStartDuty: true,
				// No error is expected for this duty because it don't overwrite (in state) the previous message
			},
			{
				Name: fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterAndSyncCommitteeDuties(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusWrongAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[3], 3, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingAttestationResponseBeaconObject(ks, version)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks, version)),
				},
				DontStartDuty: true,
				// No error is expected for this duty because it don't overwrite (in state) the previous message
			},
		}...)
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
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
			},
			PostDutyRunnerStateRoot: duplicateMsgDifferentRootsThenQuorumProposerSC(version).Root(),
			PostDutyRunnerState:     duplicateMsgDifferentRootsThenQuorumProposerSC(version).ExpectedState,
			OutputMessages:          []*types.PartialSignatureMessages{},
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
			Messages: []*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusWrongProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
			},
			PostDutyRunnerStateRoot: duplicateMsgDifferentRootsThenQuorumBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     duplicateMsgDifferentRootsThenQuorumBlindedProposerSC(version).ExpectedState,
			OutputMessages:          []*types.PartialSignatureMessages{},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBlindedBeaconBlockV(ks, version)),
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
