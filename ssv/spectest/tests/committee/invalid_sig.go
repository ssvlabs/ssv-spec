package committee

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidSig ensures that a message with invalid sig raises an error.
func InvalidSig() tests.SpecTest {

	tests := []*CommitteeSpecTest{}
	expectedError := types.SSVMessageHasInvalidSignatureErrorCode

	// Keys
	numValidators := 10
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(validatorsIndexList[0])]

	// Setup msg params
	commMsgID := testingutils.CommitteeMsgID(ks)
	version := spec.DataVersionElectra
	height := qbft.Height(testingutils.TestingDutySlotV(version))

	// Function to invalidate signature
	invalidateSig := func(msg *types.SignedSSVMessage) *types.SignedSSVMessage {
		invalidSig, err := types.SignSSVMessage(ks.OperatorKeys[3], msg.SSVMessage)
		if err != nil {
			panic(err)
		}
		msg.Signatures[0] = invalidSig
		return msg
	}

	// Duties for aggregator cases
	aggDuty := testingutils.TestingAggregatorCommitteeDuty(validatorsIndexList, []int{}, version)
	sccDuty := testingutils.TestingAggregatorCommitteeDuty([]int{}, validatorsIndexList, version)
	aggAndSccDuty := testingutils.TestingAggregatorCommitteeDuty(validatorsIndexList, validatorsIndexList, version)

	tests = append(tests, []*CommitteeSpecTest{
		{
			Name:      "attestation",
			Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
			Input: []interface{}{
				testingutils.TestingAttesterDutyForValidators(version, validatorsIndexList),
				invalidateSig(testingutils.TestingProposalMessageWithIdentifierAndFullData(
					ks.OperatorKeys[1], types.OperatorID(1), commMsgID, testingutils.TestBeaconVoteByts,
					height)),
			},
			ExpectedErrorCode: expectedError,
		},
		{
			Name:      "sync committee",
			Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
			Input: []interface{}{
				testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
				invalidateSig(testingutils.TestingProposalMessageWithIdentifierAndFullData(
					ks.OperatorKeys[1], types.OperatorID(1), commMsgID, testingutils.TestBeaconVoteByts,
					height)),
			},
			ExpectedErrorCode: expectedError,
		},
		{
			Name:      "attestations and sync committees",
			Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
			Input: []interface{}{
				testingutils.TestingCommitteeDuty(validatorsIndexList, validatorsIndexList, version),
				invalidateSig(testingutils.TestingProposalMessageWithIdentifierAndFullData(
					ks.OperatorKeys[1], types.OperatorID(1), commMsgID, testingutils.TestBeaconVoteByts,
					height)),
			},
			ExpectedErrorCode: expectedError,
		},
		{
			Name:      "aggregator",
			Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
			Input: []interface{}{
				aggDuty,

				invalidateSig(testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1)))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1),
			},
			ExpectedErrorCode: expectedError,
		},
		{
			Name:      "sync committee contributor",
			Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
			Input: []interface{}{
				sccDuty,
				invalidateSig(testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(sccDuty, ksMap, 1)))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(sccDuty, ksMap, 1),
			},
			ExpectedErrorCode: expectedError,
		},
		{
			Name:      "aggregator and sync committee contributor",
			Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
			Input: []interface{}{
				aggAndSccDuty,
				invalidateSig(testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggAndSccDuty, ksMap, 1)))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggAndSccDuty, ksMap, 1),
			},
			ExpectedErrorCode: expectedError,
		},
	}...)

	multiSpecTest := NewMultiCommitteeSpecTest(
		"invalid signature",
		testdoc.CommitteeInvalidSigDoc,
		tests,
		ks,
	)

	return multiSpecTest
}
