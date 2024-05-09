package committee

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartSyncCommittees starts a cluster runner with a list of 30 sync committee duties
func StartSyncCommittees() tests.SpecTest {

	ksMapFor30Validators := testingutils.KeySetMapForValidatorIndexList(testingutils.ValidatorIndexList(30))

	multiSpecTest := &MultiCommitteeSpecTest{
		Name: "start sync committees",
		Tests: []*CommitteeSpecTest{
			{
				Name:      "30 sync committees",
				Committee: testingutils.BaseCommittee(ksMapFor30Validators),
				Input: []interface{}{
					testingutils.TestingCommitteeSyncCommitteeDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(30)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	return multiSpecTest
}
