package committeemultipleduty

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ShuffledHappyFlowDutiesWithDifferentValidators performs the happy flow of duties with shuffled input messages (that preserves order between duty messages)
// The duties are assigned to different validators
func ShuffledHappyFlowDutiesWithDifferentValidators() tests.SpecTest {

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "shuffled happy flow duties with different validators",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		// TODO add 500
		for _, numSequencedDuties := range []int{1, 2, 4} {
			for _, numValidators := range []int{8, 30} {

				ksMap := testingutils.KeySetMapForValidators(numValidators)
				shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

				multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
					{
						Name:                   fmt.Sprintf("%v duties %v attestation (%s)", numSequencedDuties, numValidators, version.String()),
						Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
						Input:                  testingutils.CommitteeInputForDutiesWithShuffleAndDifferentValidators(numSequencedDuties, numValidators, 0, true, version),
						OutputMessages:         testingutils.CommitteeOutputMessagesForDutiesWithDifferentValidators(numSequencedDuties, numValidators, 0, version),
						BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithDifferentValidators(numSequencedDuties, numValidators, 0, version),
					},
					{
						Name:                   fmt.Sprintf("%v duties %v sync committee (%s)", numSequencedDuties, numValidators, version.String()),
						Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
						Input:                  testingutils.CommitteeInputForDutiesWithShuffleAndDifferentValidators(numSequencedDuties, 0, numValidators, true, version),
						OutputMessages:         testingutils.CommitteeOutputMessagesForDutiesWithDifferentValidators(numSequencedDuties, 0, numValidators, version),
						BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithDifferentValidators(numSequencedDuties, 0, numValidators, version),
					},
					{
						Name:                   fmt.Sprintf("%v duties %v attestations %v sync committees (%s)", numSequencedDuties, numValidators, numValidators, version.String()),
						Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
						Input:                  testingutils.CommitteeInputForDutiesWithShuffleAndDifferentValidators(numSequencedDuties, numValidators, numValidators, true, version),
						OutputMessages:         testingutils.CommitteeOutputMessagesForDutiesWithDifferentValidators(numSequencedDuties, numValidators, numValidators, version),
						BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithDifferentValidators(numSequencedDuties, numValidators, numValidators, version),
					},
				}...)
			}
		}
	}

	return multiSpecTest
}
