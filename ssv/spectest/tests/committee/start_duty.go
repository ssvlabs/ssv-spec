package committee

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartDuty starts a cluster runner
func StartDuty() tests.SpecTest {

	multiSpecTest := &MultiCommitteeSpecTest{
		Name:  "start duty",
		Tests: []*CommitteeSpecTest{},
	}

	for _, numValidators := range []int{1, 30, 100, 500} {

		validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
		ksMap := testingutils.KeySetMapForValidators(numValidators)

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*CommitteeSpecTest{

			{
				Name:      fmt.Sprintf("%v attestation", numValidators),
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingCommitteeAttesterDuty(testingutils.TestingDutySlot, validatorsIndexList),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
			{
				Name:      fmt.Sprintf("%v sync committee", numValidators),
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingCommitteeSyncCommitteeDuty(testingutils.TestingDutySlot, validatorsIndexList),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
			{
				Name:      fmt.Sprintf("%v attestation %v sync committee", numValidators, numValidators),
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, validatorsIndexList, validatorsIndexList),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		}...)
	}

	return multiSpecTest
}
