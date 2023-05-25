package consensus

import (
	"encoding/hex"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	ssz "github.com/ferranbt/fastssz"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData(ks), ks, types.BNRoleSyncCommitteeContribution),
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.SyncCommitteeContributionMsgID,
							testingutils.TestSyncCommitteeContributionConsensusDataByts(ks),
						), nil),
				),
				PostDutyRunnerStateRoot: "5448f47bb76e4639629e146c242cb27a4a265fae9d871f0fc0f3c66aeea60eb0",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee),
					testingutils.SSVMsgSyncCommittee(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.SyncCommitteeMsgID,
							testingutils.TestSyncCommitteeConsensusDataByts,
						), nil),
				),
				PostDutyRunnerStateRoot: "8097eefe5cbd4590de980ac66db0f033552f608fef580d1b9c452bcf2f1513fd",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData(ks), ks, types.BNRoleAggregator),
					testingutils.SSVMsgAggregator(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.AggregatorMsgID,
							testingutils.TestAggregatorConsensusDataByts(ks),
						), nil),
				),
				PostDutyRunnerStateRoot: "99ddd0aadc2708e0a33f9dd979fd45af9bf71d8d85fc1b04cbb0e418e909bcd4",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester),
					testingutils.SSVMsgAttester(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.AttesterMsgID,
							testingutils.TestAttesterConsensusDataByts,
						), nil),
				),
				PostDutyRunnerStateRoot: "4016e9bc1405a443f4a2755f5927d9017c66dd1f5246d03c9bcf7b353649a460",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
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
				testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.BNRoleProposer), // consensus
				[]*types.SSVMessage{ // post consensus
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version)),
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.ProposerMsgID,
							testingutils.TestProposerConsensusDataBytsV(version),
						), nil),
				}...,
			),
			PostDutyRunnerStateRoot: postFinishProposerSC(version).Root(),
			PostDutyRunnerState:     postFinishProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				getSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
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
				testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(version), ks, types.BNRoleProposer), // consensus
				[]*types.SSVMessage{ // post consensus
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)),
					testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version)),
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, testingutils.ProposerMsgID,
							testingutils.TestProposerBlindedBlockConsensusDataBytsV(version),
						), nil),
				}...,
			),
			PostDutyRunnerStateRoot: postFinishBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     postFinishBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				getSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
			},
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}

func getSSZRootNoError(obj ssz.HashRoot) string {
	r, _ := obj.HashTreeRoot()
	return hex.EncodeToString(r[:])
}
