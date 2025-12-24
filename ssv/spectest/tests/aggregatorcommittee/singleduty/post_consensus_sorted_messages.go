package singleduty

import (
	"crypto/sha256"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SortedPostConsensusMessages tests that the post-consensus messages are sorted by (validator index, signing root)
func SortedPostConsensusMessages() tests.SpecTest {

	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(1)]
	msgID := testingutils.AggregatorCommitteeMsgID(ks)

	var testCases []*committee.CommitteeSpecTest

	for _, version := range testingutils.SupportedAggregatorVersions {
		// Test with unsorted validator indices: 5, 2, 8, 1
		// After sorting, should be: 1, 2, 5, 8
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
		slot := testingutils.TestingDutySlotV(version)
		height := qbft.Height(slot)

		consensusDataBytes := testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(duty, version)

		testCases = append(testCases, &committee.CommitteeSpecTest{
			Name: fmt.Sprintf("agg (%s) and scc", version.String()),
			Committee: testingutils.
				BaseAggregatorCommitteeWithCreatorFieldsFromRunner(remappedKsMap, testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap).(*ssv.AggregatorCommitteeRunner)),
			Input: []interface{}{
				duty,

				// Pre-consensus messages
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, remappedKsMap, 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, remappedKsMap, 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, remappedKsMap, 3, version))),

				// Consensus messages
				testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1), msgID, consensusDataBytes, height),
				testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
				testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
				testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
				testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
				testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
				testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				// Pre-consensus message broadcasted when starting duty (sorted)
				testingutils.PreConsensusAggregatorCommitteeMsgForDutySorted(duty, remappedKsMap, 1, version),
				// Post-consensus message broadcasted after consensus (sorted)
				testingutils.PostConsensusAggregatorCommitteeMsgForDutySorted(duty, remappedKsMap, 1, version),
			},
			BeaconBroadcastedRoots: []string{},
			StrictMessageOrder:     true,
		})
	}

	multiSpecTest := committee.NewMultiCommitteeSpecTest(
		"aggregator committee sorted post-consensus messages",
		testdoc.AggregatorCommitteeDutyPostConsensusSortingDoc,
		testCases,
		ks,
	)

	return multiSpecTest
}
