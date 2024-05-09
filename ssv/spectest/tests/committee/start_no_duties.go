package committee

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartNoDuty starts a cluster runner with no duties
// Expects: error?
func StartNoDuty() tests.SpecTest {

	ksMapFor1Validator := testingutils.KeySetMapForValidatorIndexList(testingutils.ValidatorIndexList(1))

	multiSpecTest := &MultiCommitteeSpecTest{
		Name: "start no duties",
		Tests: []*CommitteeSpecTest{
			{
				Name:      "no duty",
				Committee: testingutils.BaseCommittee(ksMapFor1Validator),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, nil, nil),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	return multiSpecTest
}
