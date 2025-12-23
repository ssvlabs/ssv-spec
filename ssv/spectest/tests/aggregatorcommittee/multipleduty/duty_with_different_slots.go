package multipleduty

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

	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(1)]

	var testCases []*committee.CommitteeSpecTest

	for _, version := range testingutils.SupportedAttestationVersions {
		for _, numValidators := range []int{1, 30} {

			ksMap := testingutils.KeySetMapForValidators(numValidators)
			shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

			duty := testingutils.TestingAggregatorAndSyncCommitteeContributorDutiesWithDifferentSlot(version)
			testCases = append(testCases, []*committee.CommitteeSpecTest{
				{
					Name: fmt.Sprintf("%v aggregator (%s) and scc", numValidators, version.String()),
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
	}

	multiSpecTest := committee.NewMultiCommitteeSpecTest(
		"aggregator committee runner duty with different slots",
		testdoc.AggregatorCommitteeDutyWithDifferentSlotsDoc,
		testCases,
		ks,
	)

	return multiSpecTest
}
