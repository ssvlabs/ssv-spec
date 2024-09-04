package validation

import (
	"fmt"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ConsensusMessageLateSlot tests a consensus message received at a late slot (considering the duty's starting slot)
func ConsensusMessageLateSlot() tests.SpecTest {

	expectedErr := validation.ErrLateSlotMessage.Error()

	// Test cases: Role -> Slot in which the message is received (after duty's starting slot) -> Expected error
	roleTestCases := map[types.RunnerRole]map[int]string{
		types.RoleCommittee: {
			34: "",
			35: expectedErr,
		},
		types.RoleAggregator: {
			34: "",
			35: expectedErr,
		},
		types.RoleProposer: {
			3: "",
			4: expectedErr,
		},
		types.RoleSyncCommitteeContribution: {
			3: "",
			4: expectedErr,
		},
	}

	// Create multi test
	multiTest := &MultiMessageValidationTest{
		Name:  "consensus message sent late for slot",
		Tests: make([]*MessageValidationTest, 0),
	}

	// Add test cases
	bn := testingutils.NewTestingBeaconNode().GetBeaconNetwork()
	for role, slotDelayMap := range roleTestCases {
		for slotDelay, testCaseErr := range slotDelayMap {

			// Set received time considering slot delay
			receivedAt := time.Unix(bn.EstimatedTimeAtSlot(phase0.Slot(qbft.FirstHeight)), 0)
			receivedAt = receivedAt.Add(bn.SlotDurationSec() * time.Duration(slotDelay))

			multiTest.Tests = append(multiTest.Tests, &MessageValidationTest{
				Name:          fmt.Sprintf("role %v with %v slots delay", role, slotDelay),
				Messages:      [][]byte{testingutils.EncodeMessage(testingutils.ConsensusMsgForSlot(phase0.Slot(qbft.FirstHeight), testingutils.MessageIDForRole(role), testingutils.DefaultKeySet))},
				ExpectedError: testCaseErr,
				ReceivedAt:    receivedAt,
			})
		}
	}

	return multiTest
}

// PartialSignatureMessageLateSlot tests a partial signature message received at a late slot (considering the duty's starting slot)
func PartialSignatureMessageLateSlot() tests.SpecTest {

	expectedErr := validation.ErrLateSlotMessage.Error()

	// Test cases: Role -> Slot in which the message is received (after duty's starting slot) -> Expected error
	roleTestCases := map[types.RunnerRole]map[int]string{
		types.RoleCommittee: {
			34: "",
			35: expectedErr,
		},
		types.RoleAggregator: {
			34: "",
			35: expectedErr,
		},
		types.RoleProposer: {
			3: "",
			4: expectedErr,
		},
		types.RoleSyncCommitteeContribution: {
			3: "",
			4: expectedErr,
		},
		types.RoleValidatorRegistration: {
			100: "",
		},
		types.RoleVoluntaryExit: {
			100: "",
		},
	}

	// Create multi test
	multiTest := &MultiMessageValidationTest{
		Name:  "partial signature message sent late for slot",
		Tests: make([]*MessageValidationTest, 0),
	}

	// Add test cases
	bn := testingutils.NewTestingBeaconNode().GetBeaconNetwork()
	for role, slotDelayMap := range roleTestCases {
		for slotDelay, testCaseErr := range slotDelayMap {

			// Set received time considering slot delay
			receivedAt := time.Unix(bn.EstimatedTimeAtSlot(0), 0)
			receivedAt = receivedAt.Add(bn.SlotDurationSec() * time.Duration(slotDelay))

			multiTest.Tests = append(multiTest.Tests, &MessageValidationTest{
				Name:          fmt.Sprintf("role %v with %v slots delay", role, slotDelay),
				Messages:      [][]byte{testingutils.EncodeMessage(testingutils.PartialSignatureMsgForSlot(0, testingutils.MessageIDForRole(role), 1, testingutils.DefaultKeySet))},
				ExpectedError: testCaseErr,
				ReceivedAt:    receivedAt,
			})
		}
	}

	return multiTest
}
