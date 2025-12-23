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

// SortedPreConsensusMessages tests that the pre-consensus messages are sorted by (validator index, signing root)
func SortedPreConsensusMessages() tests.SpecTest {

	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(1)]

	var testCases []*committee.CommitteeSpecTest

	for _, version := range testingutils.SupportedAttestationVersions {
		unsortedValidatorIndices := []int{5, 2, 8, 1}

		ksMap := testingutils.KeySetMapForValidators(len(unsortedValidatorIndices))
		// Remap to our specific unsorted indices
		remappedKsMap := make(map[phase0.ValidatorIndex]*testingutils.TestKeySet)
		for i, valIdx := range unsortedValidatorIndices {
			originalKs := ksMap[phase0.ValidatorIndex(i+1)]
			remappedKsMap[phase0.ValidatorIndex(valIdx)] = originalKs
		}

		shareMap := testingutils.ShareMapFromKeySetMap(remappedKsMap)

		// Create duty with unsorted validator indices
		duty := testingutils.TestingAggregatorDutyForValidators(version, unsortedValidatorIndices)

		testCases = append(testCases, &committee.CommitteeSpecTest{
			Name: fmt.Sprintf("agg (%s) and scc", version.String()),
			Committee: testingutils.
				BaseAggregatorCommitteeWithCreatorFieldsFromRunner(remappedKsMap, testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap).(*ssv.AggregatorCommitteeRunner)),
			Input: []interface{}{
				duty,
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDutySorted(duty, remappedKsMap, 1, version),
			},
			BeaconBroadcastedRoots: []string{},
			StrictMessageOrder:     true,
		})
	}

	multiSpecTest := committee.NewMultiCommitteeSpecTest(
		"aggregator committee sorted pre-consensus messages",
		testdoc.AggregatorCommitteeDutyPreConsensusSortingDoc,
		testCases,
		ks,
	)

	return multiSpecTest
}
