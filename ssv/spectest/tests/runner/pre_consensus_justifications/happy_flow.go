package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow  tests a full flow of an already started duty with pre-consensus justifications
func HappyFlow() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	consensusMsgs := func(cd *types.ConsensusData, role types.BeaconRole, height qbft.Height) []*types.SSVMessage {
		id := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], role)
		qbftMsgs := testingutils.SSVDecidingMsgsForHeight(cd, id[:], height, ks)

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
		Name: "pre consensus full happy flow",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					// consensus
					consensusMsgs(testingutils.TestContributionProofWithJustificationsConsensusData(ks),
						types.BNRoleSyncCommitteeContribution, testingutils.TestingDutySlot),
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
					// consensus
					consensusMsgs(testingutils.TestSelectionProofWithJustificationsConsensusData(ks),
						types.BNRoleAggregator, testingutils.TestingDutySlot),
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
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks)),
				},
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: append(
					// consensus
					consensusMsgs(testingutils.TestProposerWithJustificationsConsensusDataV(ks,
						spec.DataVersionBellatrix), types.BNRoleProposer, testingutils.TestingDutySlot),
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, spec.DataVersionBellatrix)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, spec.DataVersionBellatrix)),
					}...,
				),
				PostDutyRunnerStateRoot: "020e7b6eb0853ffab8dab87cc6ff2ed66d7c5d2da0280288775bb040c3bec524",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, spec.DataVersionBellatrix)),
				},
			},
			{
				Name:   "proposer blinded block",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: append(
					// consensus
					consensusMsgs(testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks,
						spec.DataVersionBellatrix), types.BNRoleProposer, testingutils.TestingDutySlot),
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, spec.DataVersionBellatrix)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, spec.DataVersionBellatrix)),
					}...,
				),
				PostDutyRunnerStateRoot: "f6ee425da198f20ec945c1764e20a7c9f5b48815482add50a64b3287e700216b",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, spec.DataVersionBellatrix)),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				// we add random justifications to the consensus data.
				// this is not a valid scenario, but we want to test that the runner will ignore them.
				// This is to confirm consistent behaviors across nodes.
				Messages: append(
					consensusMsgs(testingutils.TestSyncCommitteeConsensusDataWithPreconJust(ks),
						types.BNRoleSyncCommittee,
						testingutils.TestingDutySlot),
					// consensus
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
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks)),
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: append(
					consensusMsgs(testingutils.TestAttesterConsensusDataWithPreconJust(ks), types.BNRoleAttester,
						testingutils.TestingDutySlot), // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, testingutils.TestingDutySlot)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, testingutils.TestingDutySlot)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, testingutils.TestingDutySlot)),
					}...,
				),
				PostDutyRunnerStateRoot: "9d55ff5721b21c5b99dd4b4bacb0acda0b674112fe3cec55cc6aeb04ad5dc2fc",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, testingutils.TestingDutySlot),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAttestation(ks)),
				},
			},
		},
	}
}
