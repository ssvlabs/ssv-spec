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

// HappyFlow performs a complete duty execution
func HappyFlow() tests.SpecTest {

	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(1)]
	msgID := testingutils.CommitteeMsgID(ks)

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "happy flow",
		Tests: []*committee.CommitteeSpecTest{},
	}

	// TODO add 500
	for _, numValidators := range []int{1, 30} {

		validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
		ksMap := testingutils.KeySetMapForValidators(numValidators)
		shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
			{
				Name:      fmt.Sprintf("%v attestation", numValidators),
				Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
				Input: []interface{}{
					testingutils.TestingCommitteeAttesterDuty(testingutils.TestingDutySlot, validatorsIndexList),
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
						qbft.Height(testingutils.TestingDutySlot)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 1, testingutils.TestingDutySlot))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 2, testingutils.TestingDutySlot))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 3, testingutils.TestingDutySlot))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 1, testingutils.TestingDutySlot),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedAttestationSSZRootForKeyMap(ksMap),
			},
			{
				Name:      fmt.Sprintf("%v sync committee", numValidators),
				Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
				Input: []interface{}{
					testingutils.TestingCommitteeSyncCommitteeDuty(testingutils.TestingDutySlot, validatorsIndexList),
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
						qbft.Height(testingutils.TestingDutySlot)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 3))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 1),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(ksMap),
			},
			{
				Name:      fmt.Sprintf("%v attestations %v sync committees", numValidators, numValidators),
				Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, validatorsIndexList, validatorsIndexList),
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
						qbft.Height(testingutils.TestingDutySlot)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 1, testingutils.TestingDutySlot))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 2, testingutils.TestingDutySlot))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 3, testingutils.TestingDutySlot))),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 1, testingutils.TestingDutySlot),
				},
				BeaconBroadcastedRoots: append(
					testingutils.TestingSignedAttestationSSZRootForKeyMap(ksMap),
					testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(ksMap)...),
			},
		}...)
	}

	return multiSpecTest
}
