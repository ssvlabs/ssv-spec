package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TerminatesWhenAllSignersSeen ensures the runner terminates after seeing messages from all committee signers
func TerminatesWhenAllSignersSeen() tests.SpecTest {
	// Create key share map for several validators, though only one will be used to form quorum
	ks := testingutils.Testing4SharesSet()
	numValidators := 10
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus terminates when all signers seen",
		testdoc.PreConsensusAggregatorCommitteeAllMessagesReceivedAndNoAggregatorsDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {

		// Duty with several validators to be used as input
		mixedDuty := testingutils.TestingAggregatorCommitteeDuty(validatorsIndexList, validatorsIndexList, version)
		slot := mixedDuty.Slot

		// Extract one duty for pre-consensus messages
		// Get sample aggregator duty
		idx, firstAggregatorDuty := GetFirstAggregatorDuty(mixedDuty)
		// Bump committee index to differentiate from others (for beacon node selection)
		firstAggregatorDuty.CommitteeIndex += 1
		mixedDuty.ValidatorDuties[idx] = firstAggregatorDuty // store updated duty back in mixed duty

		// Created single-duty view with aggregator duty sample
		singleValidatorDuty := &types.AggregatorCommitteeDuty{
			Slot:            slot,
			ValidatorDuties: []*types.ValidatorDuty{firstAggregatorDuty},
		}

		// Create beacon node that does NOT select it as aggregator
		beacon := testingutils.NewTestingBeaconNode()
		beacon.SetAggregators(map[phase0.CommitteeIndex]bool{
			firstAggregatorDuty.CommitteeIndex: false, // not selected as aggregator
		})

		runner, err := testingutils.ConstructBaseRunnerWithShareMapAndBeaconNode(types.RoleAggregatorCommittee, shareMap, beacon)
		if err != nil {
			panic(err)
		}

		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("terminates when all signers seen (%s)", version.String()),
			Runner: runner,
			Duty:   mixedDuty, // Send full duty to runner
			Messages: []*types.SignedSSVMessage{
				// Send, for all operators, a pre-consensus messages only for the single validator that won't be selected
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 1))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 2))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 3))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 4))),

				// Post-consensus message raises no running duty error
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 1, version))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1),
			},
			// Future message error is expected since runner should not have started consensus and thus QBFT controller hasn't been updated to such height
			BeaconAggregators:       []phase0.CommitteeIndex{firstAggregatorDuty.CommitteeIndex},
			BeaconAggregatorsValues: []bool{false},
			ExpectedErrorCode:       types.NoRunningDutyErrorCode,
		})
	}

	return multiSpecTest
}
