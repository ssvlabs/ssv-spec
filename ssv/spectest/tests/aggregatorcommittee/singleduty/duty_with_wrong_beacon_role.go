package singleduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DutyWithWrongBeaconRole tries to execute an aggregator committee duty that also holds an attester ValidatorDuty for an owned validator.
func DutyWithWrongBeaconRole() tests.SpecTest {

	ksMap := testingutils.KeySetMapForValidators(2)
	ks := ksMap[phase0.ValidatorIndex(1)]
	var testCases []*committee.CommitteeSpecTest
	for _, version := range testingutils.SupportedAggregatorVersions {
		// Aggregator duty first: JSON runners infer the duty kind from the first ValidatorDuty's role
		duty := testingutils.TestingAggregatorCommitteeDuty([]int{1}, nil, version)
		attesterDuty := testingutils.TestingAttesterDutyForValidators(version, []int{2}).ValidatorDuties[0]
		duty.ValidatorDuties = append(duty.ValidatorDuties, attesterDuty)

		testCases = append(testCases, &committee.CommitteeSpecTest{
			Name:                   fmt.Sprintf("aggregator with attester (%s)", version.String()),
			Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
			Input:                  []interface{}{duty},
			OutputMessages:         []*types.PartialSignatureMessages{},
			BeaconBroadcastedRoots: []string{},
			ExpectedErrorCode:      types.InvalidAggregatorCommitteeDutyErrorCode,
		})
	}

	multiSpecTest := committee.NewMultiCommitteeSpecTest(
		"aggregator committee runner duty with wrong beacon role",
		testdoc.AggregatorCommitteeDutyWithWrongBeaconRoleDoc,
		testCases,
		ks,
	)

	return multiSpecTest
}
