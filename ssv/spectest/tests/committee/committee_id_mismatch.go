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

// CommitteeIDMismatch ensures that a message with a different committee ID in the MessageID field raises an error.
func CommitteeIDMismatch() tests.SpecTest {

	tests := []*CommitteeSpecTest{}
	expectedError := types.MessageIDCommitteeIDMismatchErrorCode

	// Keys
	numValidators := 10
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(validatorsIndexList[0])]

	// Setup msg params
	commMsgID := testingutils.CommitteeMsgID(ks)
	version := spec.DataVersionElectra
	height := qbft.Height(testingutils.TestingDutySlotV(version))

	// Function to change committee ID
	invalidateCommitteeID := func(msg *types.SignedSSVMessage) *types.SignedSSVMessage {
		// Change committee ID to wrong one
		msgID := msg.SSVMessage.MsgID
		wrongCommitteeID := [48]byte{1}
		msg.SSVMessage.MsgID = types.NewMsgID(
			types.DomainType(msgID.GetDomain()),
			wrongCommitteeID[:],
			msgID.GetRoleType(),
		)

		// fix signature not to raise invalid sig error
		newSignature, err := types.SignSSVMessage(ks.OperatorKeys[1], msg.SSVMessage)
		if err != nil {
			panic(err)
		}
		msg.Signatures[0] = newSignature

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
				invalidateCommitteeID(testingutils.TestingProposalMessageWithIdentifierAndFullData(
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
				invalidateCommitteeID(testingutils.TestingProposalMessageWithIdentifierAndFullData(
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
				invalidateCommitteeID(testingutils.TestingProposalMessageWithIdentifierAndFullData(
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

				invalidateCommitteeID(testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1, version)))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1, version),
			},
			ExpectedErrorCode: expectedError,
		},
		{
			Name:      "sync committee contributor",
			Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
			Input: []interface{}{
				sccDuty,
				invalidateCommitteeID(testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(sccDuty, ksMap, 1, version)))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(sccDuty, ksMap, 1, version),
			},
			ExpectedErrorCode: expectedError,
		},
		{
			Name:      "aggregator and sync committee contributor",
			Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap),
			Input: []interface{}{
				aggAndSccDuty,
				invalidateCommitteeID(testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggAndSccDuty, ksMap, 1, version)))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggAndSccDuty, ksMap, 1, version),
			},
			ExpectedErrorCode: expectedError,
		},
	}...)

	multiSpecTest := NewMultiCommitteeSpecTest(
		"committee id mismatch",
		testdoc.CommitteeMismatchCommitteeIDDoc,
		tests,
		ks,
	)

	return multiSpecTest
}
