package committee

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartMixedDuties starts a cluster runner with 30 attestation and 30 sync committee duties
func StartMixedDuties() tests.SpecTest {

	ksMapFor30Validators := testingutils.KeySetMapForValidatorIndexList(testingutils.ValidatorIndexList(30))

	multiSpecTest := &MultiCommitteeSpecTest{
		Name: "start mixed duties",
		Tests: []*CommitteeSpecTest{
			{
				Name:      "30 attestations 30 sync committees",
				Committee: testingutils.BaseCommittee(ksMapFor30Validators),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(30), testingutils.ValidatorIndexList(30)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	return multiSpecTest
}
