package committeesingleduty

import (
	"crypto/sha256"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MissingSomeShares performs a complete duty execution for a runner that only has a fraction of the shares for the duty's validators
func MissingSomeShares() tests.SpecTest {

	// Message ID
	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(1)]
	msgID := testingutils.CommitteeMsgID(ks)

	// Committee's validator indexes
	committeeShareValidators := []int{1, 3, 5, 7, 9}
	// KeySet and Share map for Committee
	committeeShareKSMap := testingutils.KeySetMapForValidatorIndexList(committeeShareValidators)
	committeeShareMap := testingutils.ShareMapFromKeySetMap(committeeShareKSMap)

	// Duty's validator indexes
	dutyValidators := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "start committee duty with missing shares",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {

		slot := testingutils.TestingDutySlotV(version)
		height := qbft.Height(slot)

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{

			{
				Name:      fmt.Sprintf("attestation (%s)", version.String()),
				Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(committeeShareKSMap, testingutils.CommitteeRunnerWithShareMap(committeeShareMap).(*ssv.CommitteeRunner)),
				Input: []interface{}{
					// Duty for more validators
					testingutils.TestingAttesterDutyForValidators(version, dutyValidators),

					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
						height),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(committeeShareKSMap, 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(committeeShareKSMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(committeeShareKSMap, 3, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationMsgForKeySet(committeeShareKSMap, 1, version),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedAttestationResponseSSZRootForKeyMap(committeeShareKSMap, version),
			},
			{
				Name:      fmt.Sprintf("sync committee (%s)", version.String()),
				Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(committeeShareKSMap, testingutils.CommitteeRunnerWithShareMap(committeeShareMap).(*ssv.CommitteeRunner)),
				Input: []interface{}{
					// Duty for more validators
					testingutils.TestingSyncCommitteeDutyForValidators(version, dutyValidators),

					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
						height),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(committeeShareKSMap, 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(committeeShareKSMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(committeeShareKSMap, 3, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusSyncCommitteeMsgForKeySet(committeeShareKSMap, 1, version),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(committeeShareKSMap, version),
			},
			{
				Name:      fmt.Sprintf("attestations sync committees (%s)", version.String()),
				Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(committeeShareKSMap, testingutils.CommitteeRunnerWithShareMap(committeeShareMap).(*ssv.CommitteeRunner)),
				Input: []interface{}{
					// Duty for more validators
					testingutils.TestingCommitteeDuty(dutyValidators, dutyValidators, version),

					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
						height),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(committeeShareKSMap, 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(committeeShareKSMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(committeeShareKSMap, 3, version))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(committeeShareKSMap, 1, version),
				},
				BeaconBroadcastedRoots: append(
					testingutils.TestingSignedAttestationResponseSSZRootForKeyMap(committeeShareKSMap, version),
					testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(committeeShareKSMap, version)...),
			},
		}...)
	}

	return multiSpecTest
}
