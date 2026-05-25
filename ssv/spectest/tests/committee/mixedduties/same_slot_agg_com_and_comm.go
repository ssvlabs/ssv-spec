package committeemultipleduty

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SameSlot tests that both a committee duty and an aggregator committee duty can be performed on the same slot
func SameSlot() tests.SpecTest {

	vals := testingutils.ValidatorIndexList(10)
	ksMap := testingutils.KeySetMapForValidators(10)
	version := spec.DataVersionElectra
	slot := testingutils.TestingDutySlotV(version)

	// Committee duty and Aggregator committee duty on same slot
	aggDuty, aggMsgs := testingutils.AggregatorCommitteeInputForSlot(vals, vals, ksMap, slot, version)
	commDuty, commMsgs := testingutils.CommitteeInputForSlot(slot, vals, vals, ksMap, true)

	// Build input
	input := make([]interface{}, 0)
	input = append(input, commDuty)
	input = append(input, aggDuty)
	for _, msg := range commMsgs {
		input = append(input, msg)
	}
	for _, msg := range aggMsgs {
		input = append(input, msg)
	}

	// Build output
	outputMsgs := []*types.PartialSignatureMessages{
		// Comm
		testingutils.PostConsensusCommitteeMsgForDuty(commDuty, ksMap, 1),
		// Agg Comm
		testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1),
		testingutils.PostConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1, version),
	}

	// Build beacon roots
	beaconRoots := make([]string, 0)
	beaconRoots = append(beaconRoots, testingutils.TestingSignedCommitteeBeaconObjectSSZRoot(commDuty, ksMap, version)...)
	beaconRoots = append(beaconRoots, testingutils.TestingSignedAggregatorCommitteeBeaconObjectSSZRoot(aggDuty, ksMap, version)...)

	tests := []*committee.CommitteeSpecTest{
		{
			Name:                   "mixed duties",
			Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
			Input:                  input,
			OutputMessages:         outputMsgs,
			BeaconBroadcastedRoots: beaconRoots,
		},
	}

	multiSpecTest := committee.NewMultiCommitteeSpecTest(
		"committee with mixed committee duty and aggregator committee duty on same slot",
		testdoc.CommitteeMixedDutiesDoc,
		tests,
		nil,
	)
	return multiSpecTest
}
