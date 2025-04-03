package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicateMsgDifferentRoots tests duplicate SignedPartialSignatureMessage (from same signer) but with different roots
func DuplicateMsgDifferentRoots() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	expectedError := "failed processing post consensus message: invalid post-consensus message: wrong signing root"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus duplicate msg different roots",
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
				},
				PostDutyRunnerStateRoot: duplicateMsgDifferentRootsSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     duplicateMsgDifferentRootsSyncCommitteeContributionSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots:  []string{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
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
			},
			OutputMessages:         []*types.PartialSignatureMessages{},
			BeaconBroadcastedRoots: []string{},
			DontStartDuty:          true,
			ExpectedError:          expectedError,
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
				},
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{},
				DontStartDuty:          true,
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
				},
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{},
				DontStartDuty:          true,
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
				},
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{},
				DontStartDuty:          true,
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
			},
			PostDutyRunnerStateRoot: duplicateMsgDifferentRootsProposerSC(version).Root(),
			PostDutyRunnerState:     duplicateMsgDifferentRootsProposerSC(version).ExpectedState,
			OutputMessages:          []*types.PartialSignatureMessages{},
			BeaconBroadcastedRoots:  []string{},
			DontStartDuty:           true,
			ExpectedError:           expectedError,
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
			},
			PostDutyRunnerStateRoot: duplicateMsgDifferentRootsBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     duplicateMsgDifferentRootsBlindedProposerSC(version).ExpectedState,
			OutputMessages:          []*types.PartialSignatureMessages{},
			BeaconBroadcastedRoots:  []string{},
			DontStartDuty:           true,
			ExpectedError:           expectedError,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
