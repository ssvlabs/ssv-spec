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

	consensusMsgs := func(cd *types.ConsensusData, role types.BeaconRole) []*types.SSVMessage {
		id := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], role)
		qbftMsgs := testingutils.DecidingMsgsForHeight(cd, id[:], 2, ks)

		ret := make([]*types.SSVMessage, 0)
		for _, msg := range qbftMsgs {
			byts, _ := msg.Encode()
			ret = append(ret, &types.SSVMessage{
				MsgType: types.SSVConsensusMsgType,
				MsgID:   id,
				Data:    byts,
			})
		}
		return ret
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "full happy flow",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator first height",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					// consensus
					consensusMsgs(testingutils.TestContributionProofWithJustificationsConsensusData(ks), types.BNRoleSyncCommitteeContribution)[:7],
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
				Name:   "sync committee aggregator > first height",
				Runner: decideFirstHeight(testingutils.SyncCommitteeContributionRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					// consensus
					consensusMsgs(testingutils.TestContributionProofWithJustificationsConsensusData(ks), types.BNRoleSyncCommitteeContribution)[7:14],
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks)),
					}...,
				),
				PostDutyRunnerStateRoot: "139892acf4e22593f895344770a82b78288476889e54ebebb26d9382eba4e17a",
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
				Name:   "sync committee first height",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					consensusMsgs(testingutils.TestSyncCommitteeConsensusData, types.BNRoleSyncCommittee)[:7], // consensus
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
				Name:   "sync committee > first height",
				Runner: decideFirstHeight(testingutils.SyncCommitteeRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					consensusMsgs(testingutils.TestSyncCommitteeConsensusData, types.BNRoleSyncCommittee)[7:14], // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "0488511c9416bf29f7679d316068b2ad723d7440a5d93c8f51c95b0bb6e4047f",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks)),
				},
			},
			{
				Name:   "aggregator first height",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					// consensus
					consensusMsgs(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), types.BNRoleAggregator)[:7],
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
				Name:   "aggregator > first height",
				Runner: decideFirstHeight(testingutils.AggregatorRunner(ks)),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					// consensus
					consensusMsgs(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), types.BNRoleAggregator)[7:14],
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "dfa0c360fa5a320fa7823211eb208f38fa4c8fcd417c70f9451945b478ee5a48",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks)),
				},
			},
			{
				Name:   "proposer first height",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					// consensus
					consensusMsgs(testingutils.TestProposerWithJustificationsConsensusData(ks), types.BNRoleProposer)[:7],
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
				Name:   "proposer > first height",
				Runner: decideFirstHeight(testingutils.ProposerRunner(ks)),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					// consensus
					consensusMsgs(testingutils.TestProposerWithJustificationsConsensusData(ks), types.BNRoleProposer)[7:14],
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "335e870fa3c015b9e232af2ccb48065a38789ff11725e86841aeac5b2fe9cff2",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
				},
			},
			{
				Name:   "proposer blinded block first height",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					// consensus
					consensusMsgs(testingutils.TestProposerBlindedWithJustificationsConsensusData(ks), types.BNRoleProposer)[:7],
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
				Name:   "proposer blinded block > first height",
				Runner: decideFirstHeight(testingutils.ProposerBlindedBlockRunner(ks)),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					// consensus
					consensusMsgs(testingutils.TestProposerBlindedWithJustificationsConsensusData(ks), types.BNRoleProposer)[7:14],
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "c07ee02ae01e89270711cdff91fd3c9fe0117f6b6f1d275708f14c1e48344874",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
				},
			},
			{
				Name:   "attester first height",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: append(
					consensusMsgs(testingutils.TestAttesterConsensusData, types.BNRoleAttester)[:7], // consensus
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
				Name:   "validator registration first height",
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
