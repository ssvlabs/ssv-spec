package preconsensus

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SingleContributorValidatorSelected ensures a single-validator quorum that gets selected results in starting QBFT
func SingleContributorValidatorSelected() tests.SpecTest {
	// Create key share map for several validators, though only one will be used to form quorum
	ks := testingutils.Testing4SharesSet()
	numValidators := 10
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus single contributor validator selected",
		testdoc.PreConsensusAggregatorCommitteeContributorSingleValidatorSelectedDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {

		// Duty with several validators to be used as input
		mixedDuty := testingutils.TestingAggregatorCommitteeDuty([]int{}, validatorsIndexList, version)
		slot := mixedDuty.Slot

		// Extract one duty for pre-consensus messages
		_, firstContributorDuty := GetFirstContributorDuty(mixedDuty)

		// Created single-duty view with aggregator duty sample
		singleValidatorDuty := &types.AggregatorCommitteeDuty{
			Slot:            slot,
			ValidatorDuties: []*types.ValidatorDuty{firstContributorDuty},
		}
		consensusDataBytes := testingutils.TestAggregatorCommitteeConsensusDataBytesForDutyWithKS(singleValidatorDuty, version, ksMap)

		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("single contributor validator quorum starts consensus (%s)", version.String()),
			Runner: testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
			Duty:   mixedDuty, // Send full duty to runner
			Messages: []*types.SignedSSVMessage{
				// Send pre-consensus messages only for the single validator
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 1))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 2))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 3))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1),
			},
			QBFTProposals: [][]byte{consensusDataBytes}, // Expect proposal with consensus data for the single validator duty
		})
	}

	return multiSpecTest
}
