package preconsensus

import (
	"crypto/sha256"
	"fmt"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SingleValidatorSelected ensures a single-validator quorum that gets selected results in starting QBFT
func SingleValidatorSelected() tests.SpecTest {
	// Create key share map for several validators, though only one will be used to form quorum
	ks := testingutils.Testing4SharesSet()
	numValidators := 10
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

	// Message ID for committee
	msgID := testingutils.AggregatorCommitteeMsgID(ks)

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus single validator selected",
		testdoc.PreConsensusAggregatorCommitteeSingleValidatorSelectedDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {

		// Duty with several validators to be used as input
		mixedDuty := testingutils.TestingAggregatorCommitteeDuty(validatorsIndexList, validatorsIndexList, version)
		slot := mixedDuty.Slot
		height := qbft.Height(slot)

		// Extract one duty for pre-consensus messages
		_, firstAggregatorDuty := GetFirstAggregatorDuty(mixedDuty)

		// Created single-duty view with aggregator duty sample
		singleValidatorDuty := &types.AggregatorCommitteeDuty{
			Slot:            slot,
			ValidatorDuties: []*types.ValidatorDuty{firstAggregatorDuty},
		}
		consensusDataBytes := testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(singleValidatorDuty, version)

		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("single validator quorum starts consensus (%s)", version.String()),
			Runner: testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
			Duty:   mixedDuty, // Send full duty to runner
			Messages: []*types.SignedSSVMessage{
				// Send pre-consensus messages only for the single validator
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 1))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 2))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 3))),

				// Consensus messages
				testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1), msgID, consensusDataBytes, height),
				testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
				testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
				testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
				testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
				testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
				testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(consensusDataBytes)),

				// // Post-consensus messages
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 1, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 2, version))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 3, version))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(mixedDuty, ksMap, 1),
				testingutils.PostConsensusAggregatorCommitteeMsgForDuty(singleValidatorDuty, ksMap, 1, version),
			},
			BeaconBroadcastedRoots: testingutils.TestingSignedAggregatorCommitteeBeaconObjectSSZRoot(singleValidatorDuty, ksMap, version),
		})
	}

	return multiSpecTest
}
