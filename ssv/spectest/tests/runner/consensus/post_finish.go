package consensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostFinish tests a valid commit msg after runner finished
func PostFinish() tests.SpecTest {

	panic("implement me")

	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post finish",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					testingutils.SignedSSVMessageListF(ks, testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution)),
					// post consensus
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks))),
					// commit msg
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]},
							[]types.OperatorID{4},
							qbft.Height(testingutils.TestingDutySlot),
							testingutils.SyncCommitteeContributionMsgID,
							testingutils.TestSyncCommitteeContributionConsensusDataByts,
						), nil)),
				),
				PostDutyRunnerStateRoot: postFinishSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     postFinishSyncCommitteeContributionSC().ExpectedState,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					testingutils.SignedSSVMessageListF(ks, testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData, ks, types.BNRoleAggregator)),
					// post consensus
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3))),
					// commit msg
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(
						testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]},
							[]types.OperatorID{4},
							qbft.Height(testingutils.TestingDutySlot),
							testingutils.AggregatorMsgID,
							testingutils.TestAggregatorConsensusDataByts,
						), nil)),
				),
				PostDutyRunnerStateRoot: postFinishAggregatorSC().Root(),
				PostDutyRunnerState:     postFinishAggregatorSC().ExpectedState,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks)),
				},
			},
			{
				Name: "attester and sync committee",
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer (%s)", version.String()),
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SignedSSVMessageListF(ks, testingutils.SSVDecidingMsgsV(
					testingutils.TestProposerConsensusDataV(version),
					ks,
					types.BNRoleProposer,
				)), // consensus
				// post consensus
				testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
				testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
				// commit msg
				testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*bls.SecretKey{ks.Shares[4]},
						[]types.OperatorID{4},
						qbft.Height(testingutils.TestingDutySlotV(version)),
						testingutils.ProposerMsgID,
						testingutils.TestProposerConsensusDataBytsV(version),
					), nil)),
			),
			PostDutyRunnerStateRoot: postFinishProposerSC(version).Root(),
			PostDutyRunnerState:     postFinishProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
			},
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SignedSSVMessageListF(ks, testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(version), ks, types.BNRoleProposer)), // consensus
				// post consensus
				testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
				testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
				testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
				// commit msg
				testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*bls.SecretKey{ks.Shares[4]},
						[]types.OperatorID{4},
						qbft.Height(testingutils.TestingDutySlotV(version)),
						testingutils.ProposerMsgID,
						testingutils.TestProposerBlindedBlockConsensusDataBytsV(version),
					), nil)),
			),
			PostDutyRunnerStateRoot: postFinishBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     postFinishBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBlindedBeaconBlockV(ks, version)),
			},
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
