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

// WrongRole ensures that a message with wrong runner role (nor CommitteeRuner or AggregatorCommitteeRunner) in the MessageID field raises an error.
func WrongRole() tests.SpecTest {

	tests := []*CommitteeSpecTest{}
	expectedError := types.CommitteeWrongRoleErrorCode

	// Keys
	numValidators := 10
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(validatorsIndexList[0])]

	// Setup msg params
	commMsgID := testingutils.CommitteeMsgID(ks)
	version := spec.DataVersionElectra
	height := qbft.Height(testingutils.TestingDutySlotV(version))

	// Function to change runner role
	invalidateRunnerRole := func(msg *types.SignedSSVMessage) *types.SignedSSVMessage {
		// Change runner role to wrong one
		msgID := msg.SSVMessage.MsgID
		domain := types.DomainType(msgID.GetDomain())
		var committeeID types.CommitteeID
		executorID := msgID.GetDutyExecutorID()
		copy(committeeID[:], executorID[len(executorID)-len(committeeID):])
		wrongRunnerRole := types.RunnerRole(99) // Assuming 99 is an invalid role
		msg.SSVMessage.MsgID = types.NewCommitteeMsgID(domain, committeeID, wrongRunnerRole)

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
				invalidateRunnerRole(testingutils.TestingProposalMessageWithIdentifierAndFullData(
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
				invalidateRunnerRole(testingutils.TestingProposalMessageWithIdentifierAndFullData(
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
				invalidateRunnerRole(testingutils.TestingProposalMessageWithIdentifierAndFullData(
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

				invalidateRunnerRole(testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggDuty, ksMap, 1)))),
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
				invalidateRunnerRole(testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(sccDuty, ksMap, 1)))),
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
				invalidateRunnerRole(testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggAndSccDuty, ksMap, 1)))),
			},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusAggregatorCommitteeMsgForDuty(aggAndSccDuty, ksMap, 1),
			},
			ExpectedErrorCode: expectedError,
		},
	}...)

	multiSpecTest := NewMultiCommitteeSpecTest(
		"wrong role",
		testdoc.CommitteeWrongRoleDoc,
		tests,
		ks,
	)

	return multiSpecTest
}
