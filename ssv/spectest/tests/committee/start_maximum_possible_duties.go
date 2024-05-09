package committee

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartMaximumPossibleDuties starts a cluster runner with 500 attestation and 500 sync committee duties
func StartMaximumPossibleDuties() tests.SpecTest {

	ksMapFor500Validators := testingutils.KeySetMapForValidatorIndexList(testingutils.ValidatorIndexList(500))

	multiSpecTest := &MultiCommitteeSpecTest{
		Name: "start maximum possible duties",
		Tests: []*CommitteeSpecTest{
			{
				Name:      "500 attestations 500 sync committees",
				Committee: testingutils.BaseCommittee(ksMapFor500Validators),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(500), testingutils.ValidatorIndexList(500)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	return multiSpecTest
}
