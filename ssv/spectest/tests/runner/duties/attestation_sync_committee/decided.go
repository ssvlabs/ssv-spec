package attestationsynccommittee

import (
	"crypto/sha256"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Decided decides a maximal committee runner
func Decided() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	msgID := testingutils.CommitteeMsgID

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "500 attestations 500 sync committees",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, ValidatorIndexList(500), ValidatorIndexList(500)),
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestBeaconVoteByts,
						qbft.Height(testingutils.TestingDutySlot)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),

					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
					// TODO add decided and expected post-consensus message
					// testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, testingutils.TestingDutySlot, msgID, sha256.Sum256(testingutils.TestBeaconVoteByts)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	return multiSpecTest
}
