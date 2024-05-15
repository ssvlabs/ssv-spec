package committeesingleduty

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidConsensusData sends a proposal message to a maximal committee runner with a valid BeaconVote consensus data
func ValidConsensusData() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	ksMapFor500Validators := testingutils.KeySetMapForValidators(500)

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name: "valid consensus data",
		Tests: []*committee.CommitteeSpecTest{
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
