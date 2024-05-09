package committee

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartAttestations starts a cluster runner with a list of 30 attestation duties
func StartAttestations() tests.SpecTest {

	ksMapFor30Validators := testingutils.KeySetMapForValidatorIndexList(testingutils.ValidatorIndexList(30))

	multiSpecTest := &MultiCommitteeSpecTest{
		Name: "start attestations",
		Tests: []*CommitteeSpecTest{
			{
				Name:      "30 attestations",
				Committee: testingutils.BaseCommittee(ksMapFor30Validators),
				Input: []interface{}{
					testingutils.TestingCommitteeAttesterDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(30)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	return multiSpecTest
}
