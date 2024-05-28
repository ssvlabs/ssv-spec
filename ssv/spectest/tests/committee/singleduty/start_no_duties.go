package committeesingleduty

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartNoDuty starts a cluster runner with no duties
func StartNoDuty() tests.SpecTest {

	ksMapFor1Validator := testingutils.KeySetMapForValidators(1)

	return &committee.CommitteeSpecTest{
		Name:      "empty committee duty",
		Committee: testingutils.BaseCommittee(ksMapFor1Validator),
		Input: []interface{}{
			testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, nil, nil),
		},
		ExpectedError:  "no beacon duties",
		OutputMessages: []*types.PartialSignatureMessages{},
	}
}
