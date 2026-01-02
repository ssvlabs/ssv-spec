package singleduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DutyWithDifferentSlots tries to execute a duty with ValidatorDuty objects for different slots.
func DutyWithDifferentSlots() tests.SpecTest {

	valIdx := []int{1}
	ksMap := testingutils.KeySetMapForValidators(1)
	ks := ksMap[phase0.ValidatorIndex(1)]
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

	var testCases []*committee.CommitteeSpecTest
	for _, version := range testingutils.SupportedAggregatorVersions {
		duty := testingutils.TestingAggregatorAndSyncCommitteeContributorDutiesWithDifferentSlot(version, valIdx)
		testCases = append(testCases, []*committee.CommitteeSpecTest{
			{
				Name: fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
				Committee: testingutils.
					BaseAggregatorCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap).(*ssv.AggregatorCommitteeRunner)),
				Input: []interface{}{
					duty,
				},
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{},
				ExpectedErrorCode:      types.InvalidAggregatorCommitteeDutyErrorCode,
			},
		}...)
	}

	multiSpecTest := committee.NewMultiCommitteeSpecTest(
		"aggregator committee runner duty with different slots",
		testdoc.AggregatorCommitteeDutyWithDifferentSlotsDoc,
		testCases,
		ks,
	)

	return multiSpecTest
}
