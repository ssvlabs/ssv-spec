package committeesingleduty

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartDuty starts a cluster runner
func StartDuty() tests.SpecTest {

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "start duty",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		// TODO add 500
		for _, numValidators := range []int{1, 30} {

			validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
			ksMap := testingutils.KeySetMapForValidators(numValidators)

			multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{

				{
					Name:      fmt.Sprintf("%v attestation (%s)", numValidators, version.String()),
					Committee: testingutils.BaseCommittee(ksMap),
					Input: []interface{}{
						testingutils.TestingAttesterDutyForValidators(version, validatorsIndexList),
					},
					OutputMessages: []*types.PartialSignatureMessages{},
				},
				{
					Name:      fmt.Sprintf("%v sync committee (%s)", numValidators, version.String()),
					Committee: testingutils.BaseCommittee(ksMap),
					Input: []interface{}{
						testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
					},
					OutputMessages: []*types.PartialSignatureMessages{},
				},
				{
					Name:      fmt.Sprintf("%v attestation %v sync committee (%s)", numValidators, numValidators, version.String()),
					Committee: testingutils.BaseCommittee(ksMap),
					Input: []interface{}{
						testingutils.TestingCommitteeDuty(validatorsIndexList, validatorsIndexList, version),
					},
					OutputMessages: []*types.PartialSignatureMessages{},
				},
			}...)
		}
	}

	return multiSpecTest
}
