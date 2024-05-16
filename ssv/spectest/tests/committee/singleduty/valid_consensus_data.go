package committeesingleduty

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidConsensusData sends a proposal message to a committee runner with a valid BeaconVote consensus data
func ValidConsensusData() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	ksMapFor30Validators := testingutils.KeySetMapForValidators(30)

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name: "valid consensus data",
		Tests: []*committee.CommitteeSpecTest{
			{
				Name:      "30 attestations 30 sync committees",
				Committee: testingutils.BaseCommittee(ksMapFor30Validators),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(30), testingutils.ValidatorIndexList(30)),
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
