package committeesingleduty

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongBeaconVote sends a proposal message to a cluster runner with an invalid beacon vote
func WrongBeaconVote() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	ksMapFor30Validators := testingutils.KeySetMapForValidators(30)

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name: "wrong beacon vote",
		Tests: []*committee.CommitteeSpecTest{
			{
				Name:      "30 attestations 30 sync committees",
				Committee: testingutils.BaseCommittee(ksMapFor30Validators),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(30), testingutils.ValidatorIndexList(30)),
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
