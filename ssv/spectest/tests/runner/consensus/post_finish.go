package consensus

import (
	"crypto/rsa"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostFinish tests a valid commit msg after runner finished
func PostFinish() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post finish",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.RoleSyncCommitteeContribution),
					// post consensus
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks))),
					// commit msg
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{ks.OperatorKeys[4]},
						[]types.OperatorID{4},
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.SyncCommitteeContributionMsgID,
						testingutils.TestSyncCommitteeContributionConsensusDataByts,
					),
				),
				PostDutyRunnerStateRoot: postFinishSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     postFinishSyncCommitteeContributionSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
				ExpectedError: "failed processing consensus message: not processing consensus message since instance is already decided",
			},
		},
	}

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("aggregator (%s)", version.String()),
			Runner: testingutils.AggregatorRunner(ks),
			Duty:   testingutils.TestingAggregatorDuty(version),
			Messages: append(
				testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData(version), ks, types.RoleAggregator),
				// post consensus
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3, version))),
				// commit msg
				testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
					[]*rsa.PrivateKey{ks.OperatorKeys[4]},
					[]types.OperatorID{4},
					qbft.Height(testingutils.TestingDutySlotV(version)),
					testingutils.AggregatorMsgID,
					testingutils.TestAggregatorConsensusDataByts(version),
				),
			),
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version),
				testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks, version)),
			},
			ExpectedError: "failed processing consensus message: not processing consensus message since instance is already decided",
		},
		)
	}

	for _, version := range testingutils.SupportedAttestationVersions {

		height := qbft.Height(testingutils.TestingDutySlotV(version))

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name:   fmt.Sprintf("attester (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterDuty(version),
				Messages: append(
					testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, height),
					// post consensus
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, version))),
					// commit msg
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData([]*rsa.PrivateKey{ks.OperatorKeys[4]}, []types.OperatorID{4}, height, testingutils.CommitteeMsgID(ks), testingutils.TestBeaconVoteByts),
				),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, version),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingAttestationResponseBeaconObject(ks, version)),
				},
				ExpectedError: "failed processing consensus message: not processing consensus message since instance is already decided",
			},
			{
				Name:   fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty(version),
				Messages: append(
					testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, height),
					// post consensus
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3, version))),
					// commit msg
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData([]*rsa.PrivateKey{ks.OperatorKeys[4]}, []types.OperatorID{4}, height, testingutils.CommitteeMsgID(ks), testingutils.TestBeaconVoteByts),
				),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1, version),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks, version)),
				},
				ExpectedError: "failed processing consensus message: not processing consensus message since instance is already decided",
			},
			{
				Name:   fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages: append(
					testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, height),
					// post consensus
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[3], 3, version))),
					// commit msg
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData([]*rsa.PrivateKey{ks.OperatorKeys[4]}, []types.OperatorID{4}, height, testingutils.CommitteeMsgID(ks), testingutils.TestBeaconVoteByts),
				),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, version),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingAttestationResponseBeaconObject(ks, version)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks, version)),
				},
				ExpectedError: "failed processing consensus message: not processing consensus message since instance is already decided",
			},
		}...)
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer (%s)", version.String()),
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SSVDecidingMsgsV(
					testingutils.TestProposerConsensusDataV(version),
					ks,
					types.RoleProposer,
				), // consensus
				// post consensus
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
				// commit msg
				testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
					[]*rsa.PrivateKey{ks.OperatorKeys[4]},
					[]types.OperatorID{4},
					qbft.Height(testingutils.TestingDutySlotV(version)),
					testingutils.ProposerMsgID,
					testingutils.TestProposerConsensusDataBytsV(version),
				),
			),
			PostDutyRunnerStateRoot: postFinishProposerSC(version).Root(),
			PostDutyRunnerState:     postFinishProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
			},
			ExpectedError: "failed processing consensus message: not processing consensus message since instance is already decided",
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(version), ks, types.RoleProposer), // consensus
				// post consensus
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
				// commit msg
				testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
					[]*rsa.PrivateKey{ks.OperatorKeys[4]},
					[]types.OperatorID{4},
					qbft.Height(testingutils.TestingDutySlotV(version)),
					testingutils.ProposerMsgID,
					testingutils.TestProposerBlindedBlockConsensusDataBytsV(version),
				),
			),
			PostDutyRunnerStateRoot: postFinishBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     postFinishBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBlindedBeaconBlockV(ks, version)),
			},
			ExpectedError: "failed processing consensus message: not processing consensus message since instance is already decided",
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
