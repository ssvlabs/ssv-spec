package attestationsynccommittee

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongConsensusData sends a proposal message to a cluster runner with an invalid consensus data that can't be decoded to BeaconVote
// Expects: error
func WrongConsensusData() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "wrong consensus data",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "500 attestations 500 sync committees",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, ValidatorIndexList(500), ValidatorIndexList(500)),
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), testingutils.CommitteeMsgID, testingutils.TestWrongBeaconVoteByts,
						qbft.Height(testingutils.TestingDutySlot)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  "failed processing consensus message: could not process msg: invalid signed message: proposal not justified: proposal fullData invalid: attestation data source > target",
			},
		},
	}

	return multiSpecTest
}
