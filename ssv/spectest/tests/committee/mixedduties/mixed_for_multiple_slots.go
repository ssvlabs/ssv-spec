package committeemultipleduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MixedMultipleSlots tests that both a committee duty and an aggregator committee duty can be performed on the same slot, through multiple slots
func MixedMultipleSlots() tests.SpecTest {

	vals := testingutils.ValidatorIndexList(10)
	ksMap := testingutils.KeySetMapForValidators(10)
	version := spec.DataVersionElectra
	slot := testingutils.TestingDutySlotV(version)

	slotsAmount := 10

	input := make([]interface{}, 0)
	outputMsgs := make([]*types.PartialSignatureMessages, 0)
	beaconRoots := make([]string, 0)

	for i := 0; i < slotsAmount; i++ {

		currentSlot := slot + phase0.Slot(i)

		// Committee duty and Aggregator committee duty on same slot
		aggDuty, aggMsgs := testingutils.AggregatorCommitteeInputForSlot(vals, vals, ksMap, currentSlot, version)
		commDuty, commMsgs := testingutils.CommitteeInputForSlot(currentSlot, vals, vals, ksMap, true)
		// Build input
		input = append(input, commDuty)
		input = append(input, aggDuty)
		for _, msg := range commMsgs {
			input = append(input, msg)
		}
		for _, msg := range aggMsgs {
			input = append(input, msg)
		}

		// Build output
		outputMsgs = append(outputMsgs, []*types.PartialSignatureMessages{
			// Comm
			testingutils.PostConsensusCommitteeMsgForDuty(commDuty, ksMap, 1),
			// Agg Comm
			testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1, version),
			testingutils.PostConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1, version),
		}...)

		// Build beacon roots
		beaconRoots = append(beaconRoots, testingutils.TestingSignedCommitteeBeaconObjectSSZRoot(commDuty, ksMap, version)...)
		beaconRoots = append(beaconRoots, testingutils.TestingSignedAggregatorCommitteeBeaconObjectSSZRoot(aggDuty, ksMap, version)...)
	}

	tests := []*committee.CommitteeSpecTest{
		{
			Name:                   fmt.Sprintf("mixed duties for %d slots", slotsAmount),
			Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
			Input:                  input,
			OutputMessages:         outputMsgs,
			BeaconBroadcastedRoots: beaconRoots,
		},
	}

	multiSpecTest := committee.NewMultiCommitteeSpecTest(
		"committee with mixed committee duty and aggregator committee duty on same slot for multiple slots",
		testdoc.CommitteeMixedDutiesMultipleSlotsDoc,
		tests,
		nil,
	)
	return multiSpecTest
}
