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

// Decided decides a committee runner
func Decided() tests.SpecTest {

	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(1)]
	msgID := testingutils.CommitteeMsgID(ks)

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "decided",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		// TODO add 500
		for _, numValidators := range []int{1, 30} {

			validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
			ksMap := testingutils.KeySetMapForValidators(numValidators)
			shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

			slot := testingutils.TestingDutySlotV(version)
			height := qbft.Height(slot)

			multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
				{
					Name:      fmt.Sprintf("%v attestation (%s)", numValidators, version.String()),
					Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input: []interface{}{
						testingutils.TestingAttesterDutyForValidators(version, validatorsIndexList),
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
							height),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					},
					OutputMessages: []*types.PartialSignatureMessages{
						testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 1, version),
					},
				},
				{
					Name:      fmt.Sprintf("%v sync committee (%s)", numValidators, version.String()),
					Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input: []interface{}{
						testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
							height),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					},
					OutputMessages: []*types.PartialSignatureMessages{
						testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 1, version),
					},
				},
				{
					Name:      fmt.Sprintf("%v attestations %v sync committees (%s)", numValidators, numValidators, version.String()),
					Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input: []interface{}{
						testingutils.TestingCommitteeDuty(validatorsIndexList, validatorsIndexList, version),
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
							height),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					},
					OutputMessages: []*types.PartialSignatureMessages{
						testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 1, version),
					},
				},
			}...)
		}
	}

	return multiSpecTest
}
