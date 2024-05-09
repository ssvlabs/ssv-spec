package committee

import (
	"crypto/sha256"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Decided decides a committee runner
func Decided() tests.SpecTest {

	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(1)]

	ksMapFor1Validator := testingutils.KeySetMapForValidatorIndexList(testingutils.ValidatorIndexList(1))
	shareMapFor1Validator := testingutils.ShareMapFromKeySetMap(ksMapFor1Validator)
	ksMapFor500Validators := testingutils.KeySetMapForValidatorIndexList(testingutils.ValidatorIndexList(500))
	shareMapFor500Validators := testingutils.ShareMapFromKeySetMap(ksMapFor500Validators)

	msgID := testingutils.CommitteeMsgID(ks)

	multiSpecTest := &MultiCommitteeSpecTest{
		Name: "decided",
		Tests: []*CommitteeSpecTest{
			{
				Name:      "1 attestation",
				Committee: testingutils.BaseCommitteeWithRunnerSample(ksMapFor1Validator, testingutils.CommitteeRunnerWithShareMap(shareMapFor1Validator).(*ssv.CommitteeRunner)),
				Input: []interface{}{
					testingutils.TestingCommitteeAttesterDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(1)),
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
						qbft.Height(testingutils.TestingDutySlot)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, testingutils.TestingDutySlot),
				},
			},
			{
				Name:      "500 attestations 500 sync committees",
				Committee: testingutils.BaseCommitteeWithRunnerSample(ksMapFor500Validators, testingutils.CommitteeRunnerWithShareMap(shareMapFor500Validators).(*ssv.CommitteeRunner)),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(500), testingutils.ValidatorIndexList(500)),
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
						qbft.Height(testingutils.TestingDutySlot)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
				},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMapFor500Validators, 1, testingutils.TestingDutySlot),
				},
			},
		},
	}

	return multiSpecTest
}
