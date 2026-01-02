package singleduty

import (
	"crypto/sha256"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MaxValidators performs a complete duty execution for the maximum number of validators in an aggregator committee duty
func MaxValidators() tests.SpecTest {

	ksMap := testingutils.KeySetMapForValidators(3000)
	validatorsIndexList := testingutils.ValidatorIndexList(3000)
	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(validatorsIndexList[0])]
	msgID := testingutils.AggregatorCommitteeMsgID(ks)

	var testCases []*committee.CommitteeSpecTest

	// Add aggregator and sync committee contribution test cases
	for _, version := range testingutils.SupportedAggregatorVersions {

		slot := testingutils.TestingDutySlotV(version)
		duty := testingutils.TestingMaximumAggregatorCommitteeDutyWithParams(slot, validatorsIndexList, validatorsIndexList)
		height := qbft.Height(slot)

		consensusDataBytes := testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(duty, version)

		testCases = append(testCases, []*committee.CommitteeSpecTest{
			{
				Name: fmt.Sprintf("aggregator and sync committee contribution (%s)", version.String()),
				Committee: testingutils.
					BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
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

	multiSpecTest := committee.NewMultiCommitteeSpecTest(
		"aggregator committee runner max validators",
		testdoc.AggregatorCommitteeDutyMaxValidatorsDoc,
		testCases,
		ks,
	)

	return multiSpecTest
}
