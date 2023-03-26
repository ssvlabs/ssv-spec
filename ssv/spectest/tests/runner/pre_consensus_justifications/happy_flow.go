package pre_consensus_justifications

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	ssz "github.com/ferranbt/fastssz"
)

func getSSZRootNoError(obj ssz.HashRoot) string {
	r, _ := obj.HashTreeRoot()
	return hex.EncodeToString(r[:])
}

func removePreConsensusMsgs(msgs []*types.SSVMessage, ks *testingutils.TestKeySet) []*types.SSVMessage {
	return msgs[ks.Threshold:]
}

// HappyFlow  tests a full flow of an already started duty with pre-consensus justifications
func HappyFlow() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "full happy flow",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					// consensus
					removePreConsensusMsgs(
						testingutils.SSVDecidingMsgs(testingutils.TestContributionProofWithJustificationsConsensusData(ks), ks, types.BNRoleSyncCommitteeContribution),
						ks,
					),
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks)),
					}...,
				),
				PostDutyRunnerStateRoot: "464b4955c1b834e88e5f5599b3d904a396d61cbb5f65d8a43c39f27f5bb5274e",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee), // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "48c73f57659b69131467ef133ccb35d7de2fe96438d30bfa2b5ea63b19ead011",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks)),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					// consensus
					removePreConsensusMsgs(
						testingutils.SSVDecidingMsgs(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), ks, types.BNRoleAggregator),
						ks,
					),
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "87acfc67ebcfddbb5821bf2939c01666e2741605148ef7a403696d1a47df7d9a",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks)),
				},
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					// consensus
					removePreConsensusMsgs(
						testingutils.SSVDecidingMsgs(testingutils.TestProposerWithJustificationsConsensusData(ks), ks, types.BNRoleProposer),
						ks,
					),
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "020e7b6eb0853ffab8dab87cc6ff2ed66d7c5d2da0280288775bb040c3bec524",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
				},
			},
			{
				Name:   "proposer blinded block",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					// consensus
					removePreConsensusMsgs(
						testingutils.SSVDecidingMsgs(testingutils.TestProposerBlindedWithJustificationsConsensusData(ks), ks, types.BNRoleProposer),
						ks,
					),
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "f6ee425da198f20ec945c1764e20a7c9f5b48815482add50a64b3287e700216b",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester), // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)),
					}...,
				),
				PostDutyRunnerStateRoot: "9d55ff5721b21c5b99dd4b4bacb0acda0b674112fe3cec55cc6aeb04ad5dc2fc",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedAttestation(ks)),
				},
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3)),
				},
				PostDutyRunnerStateRoot: "f36c8b537afaba0894dbc8c87cb94466d8ac2623e9283f1c584e3d544b5f2b88",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
		},
	}
}
