package singleduty

import (
	"crypto/sha256"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// HappyFlow performs a complete duty execution for aggregator committee
func HappyFlow() tests.SpecTest {

	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(1)]
	msgID := testingutils.AggregatorCommitteeMsgID(ks)

	var testCases []*committee.CommitteeSpecTest

	// Add aggregator test cases
	for _, version := range testingutils.SupportedAttestationVersions {
		for _, numValidators := range []int{1, 30} {

			validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
			ksMap := testingutils.KeySetMapForValidators(numValidators)
			shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

			duty := testingutils.TestingAggregatorDutyForValidators(version, validatorsIndexList)
			slot := testingutils.TestingDutySlotV(version)
			height := qbft.Height(slot)

			consensusDataBytes := testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(duty, version)

			testCases = append(testCases, []*committee.CommitteeSpecTest{
				{
					Name: fmt.Sprintf("%v aggregator (%s)", numValidators, version.String()),
					Committee: testingutils.
						BaseAggregatorCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap).(*ssv.AggregatorCommitteeRunner)),
					Input: []interface{}{
						duty,

						// Pre-consensus messages
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 1, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 2, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 3, version))),

						// Consensus messages
						testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1), msgID, consensusDataBytes, height),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
						testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
						testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(consensusDataBytes)),

						// Post-consensus messages
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 1, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 2, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 3, version))),
					},
					OutputMessages: []*types.PartialSignatureMessages{
						// Pre-consensus message broadcasted when starting duty
						testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 1, version),
						// Post-consensus message broadcasted after consensus
						testingutils.PostConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 1, version),
					},
					BeaconBroadcastedRoots: testingutils.TestingSignedAggregatorCommitteeBeaconObjectSSZRoot(duty, ksMap, version),
				},
			}...)
		}
	}

	//// Add sync committee contribution test cases
	//for _, version := range testingutils.SupportedAttestationVersions {
	//	for _, numValidators := range []int{1} {
	//
	//		validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	//		ksMap := testingutils.KeySetMapForValidators(numValidators)
	//		shareMap := testingutils.ShareMapFromKeySetMap(ksMap)
	//
	//		duty := testingutils.TestingSyncCommitteeContributorDutyForValidators(version, validatorsIndexList)
	//		slot := testingutils.TestingDutySlotV(version)
	//		height := qbft.Height(slot)
	//
	//		consensusDataBytes := testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(duty, version)
	//
	//		testCases = append(testCases, []*committee.CommitteeSpecTest{
	//			{
	//				Name: fmt.Sprintf("%v sync committee contribution (%s)", numValidators, version.String()),
	//				Committee: testingutils.
	//					BaseAggregatorCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap).(*ssv.AggregatorCommitteeRunner)),
	//				Input: []interface{}{
	//					duty,
	//
	//					// Pre-consensus messages
	//					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 1, version))),
	//					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 2, version))),
	//					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 3, version))),
	//
	//					// Consensus messages
	//					testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1), msgID, consensusDataBytes, height),
	//					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
	//					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
	//					testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
	//					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
	//					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
	//					testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID, sha256.Sum256(consensusDataBytes)),
	//
	//					// Post-consensus messages
	//					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 1, version))),
	//					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 2, version))),
	//					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 3, version))),
	//				},
	//				OutputMessages: []*types.PartialSignatureMessages{
	//					// Pre-consensus message broadcasted when starting duty
	//					testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 1, version),
	//					// Post-consensus message broadcasted after consensus
	//					testingutils.PostConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 1, version),
	//				},
	//				BeaconBroadcastedRoots: testingutils.TestingSignedAggregatorCommitteeBeaconObjectSSZRoot(duty, ksMap, version),
	//			},
	//		}...)
	//	}
	//}

	multiSpecTest := committee.NewMultiCommitteeSpecTest(
		"aggregator committee runner happy flow",
		"Testing aggregator committee runner with complete duty flow for both aggregator and sync committee contribution",
		testCases,
	)

	return multiSpecTest
}
