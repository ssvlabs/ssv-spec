package committee

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongBeaconVote sends a proposal message to a cluster runner with an invalid beacon vote
func WrongBeaconVote() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	ksMapFor500Validators := testingutils.KeySetMapForValidators(500)

	multiSpecTest := &MultiCommitteeSpecTest{
		Name: "wrong beacon vote",
		Tests: []*CommitteeSpecTest{
			{
				Name:      "500 attestations 500 sync committees",
				Committee: testingutils.BaseCommittee(ksMapFor500Validators),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(500), testingutils.ValidatorIndexList(500)),
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), testingutils.CommitteeMsgID(ks), testingutils.TestWrongBeaconVoteByts,
						qbft.Height(testingutils.TestingDutySlot)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  "failed processing consensus message: could not process msg: invalid signed message: proposal not justified: proposal fullData invalid: attestation data source >= target",
			},
		},
	}

	return multiSpecTest
}
