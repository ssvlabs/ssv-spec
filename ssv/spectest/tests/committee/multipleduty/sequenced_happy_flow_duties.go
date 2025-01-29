package committeemultipleduty

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SequencedHappyFlowDuties performs the happy flow of a sequence of duties
func SequencedHappyFlowDuties() tests.SpecTest {

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "sequenced happy flow duties",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		// TODO add 500
		for _, numSequencedDuties := range []int{1, 2, 4} {
			for _, numValidators := range []int{1, 30} {

				ksMap := testingutils.KeySetMapForValidators(numValidators)
				shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

				multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
					{
						Name:                   fmt.Sprintf("%v duties %v attestation (%s)", numSequencedDuties, numValidators, version.String()),
						Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
						Input:                  testingutils.CommitteeInputForDuties(numSequencedDuties, numValidators, 0, true, version),
						OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numSequencedDuties, numValidators, 0, version),
						BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDuties(numSequencedDuties, numValidators, 0, version),
					},
					{
						Name:                   fmt.Sprintf("%v duties %v sync committee (%s)", numSequencedDuties, numValidators, version.String()),
						Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
						Input:                  testingutils.CommitteeInputForDuties(numSequencedDuties, 0, numValidators, true, version),
						OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numSequencedDuties, 0, numValidators, version),
						BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDuties(numSequencedDuties, 0, numValidators, version),
					},
					{
						Name:                   fmt.Sprintf("%v duties %v attestations %v sync committees (%s)", numSequencedDuties, numValidators, numValidators, version.String()),
						Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
						Input:                  testingutils.CommitteeInputForDuties(numSequencedDuties, numValidators, numValidators, true, version),
						OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numSequencedDuties, numValidators, numValidators, version),
						BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDuties(numSequencedDuties, numValidators, numValidators, version),
					},
				}...)
			}
		}
	}

	return multiSpecTest
}
