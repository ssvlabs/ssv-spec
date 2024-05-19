package committee

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidBeaconVote sends a proposal message to a maximal committee runner with a valid BeaconVote
func ValidBeaconVote() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	ksMapFor500Validators := testingutils.KeySetMapForValidators(500)

	multiSpecTest := &MultiCommitteeSpecTest{
		Name: "valid beacon vote",
		Tests: []*CommitteeSpecTest{
			{
				Name:      "500 attestations 500 sync committees",
				Committee: testingutils.BaseCommittee(ksMapFor500Validators),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(500), testingutils.ValidatorIndexList(500)),
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), testingutils.CommitteeMsgID(ks), testingutils.TestBeaconVoteByts,
						qbft.Height(testingutils.TestingDutySlot)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	return multiSpecTest
}
