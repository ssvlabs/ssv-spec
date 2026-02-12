package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SingleValidatorNotSelected ensures a single-validator quorum that is not selected that does not result in starting QBFT or finishing
func SingleValidatorNotSelected() tests.SpecTest {
	// Create key share map for several validators, though only one will be used to form quorum
	ks := testingutils.Testing4SharesSet()
	numValidators := 10
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

	// Message ID for committee
	msgID := testingutils.AggregatorCommitteeMsgID(ks)

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus single validator not selected",
		testdoc.PreConsensusAggregatorCommitteeSingleValidatorNotSelectedDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {

		// Duty with several validators to be used as input
		mixedDuty := testingutils.TestingAggregatorCommitteeDuty(validatorsIndexList, validatorsIndexList, version)
		slot := mixedDuty.Slot
		height := qbft.Height(slot)

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
		consensusDataBytes := testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(singleValidatorDuty, version)

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
			Name:   fmt.Sprintf("single validator quorum not selected doesn't start consensus (%s)", version.String()),
			Runner: runner,
			Duty:   mixedDuty, // Send full duty to runner
			Messages: []*types.SignedSSVMessage{
				// Send pre-consensus messages only for the single validator
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 1))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 2))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 3))),

				// Send consensus messages and expected error since runner should not have started consensus
				testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1), msgID, consensusDataBytes, height),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1),
			},
			// Future message error is expected since runner should not have started consensus and thus QBFT controller hasn't been updated to such height
			ExpectedErrorCode:       types.FutureMessageErrorCode,
			BeaconAggregators:       []phase0.CommitteeIndex{firstAggregatorDuty.CommitteeIndex},
			BeaconAggregatorsValues: []bool{false},
		})
	}

	return multiSpecTest
}

func GetFirstAggregatorDuty(aggDuty *types.AggregatorCommitteeDuty) (int, *types.ValidatorDuty) {
	return GetFirstDutyOfRole(aggDuty, types.BNRoleAggregator)
}

func GetFirstContributorDuty(aggDuty *types.AggregatorCommitteeDuty) (int, *types.ValidatorDuty) {
	return GetFirstDutyOfRole(aggDuty, types.BNRoleSyncCommitteeContribution)
}

func GetFirstDutyOfRole(aggDuty *types.AggregatorCommitteeDuty, role types.BeaconRole) (int, *types.ValidatorDuty) {
	for idx, d := range aggDuty.ValidatorDuties {
		if d.Type == role {
			return idx, d
		}
	}
	panic("no duty found for role")
}
