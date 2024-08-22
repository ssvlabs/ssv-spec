package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SlotAlreadyAdvancedForConsensusMessage tests two consensus messages with the second one with a past slot
func SlotAlreadyAdvancedForConsensusMessage() tests.SpecTest {

	msgs := [][]byte{
		testingutils.EncodeMessage(testingutils.ConsensusMsgForSlot(2, testingutils.DefaultMsgID, testingutils.DefaultKeySet)),
		testingutils.EncodeMessage(testingutils.ConsensusMsgForSlot(1, testingutils.DefaultMsgID, testingutils.DefaultKeySet)),
	}

	return &MessageValidationTest{
		Name:          "slot already advanced for consensus message",
		Messages:      msgs,
		ExpectedError: validation.ErrSlotAlreadyAdvanced.Error(),
	}
}

// SlotAlreadyAdvancedForPartialSignatureMessage tests two partial signature messages with the second one with a past slot
func SlotAlreadyAdvancedForPartialSignatureMessage() tests.SpecTest {

	msgs := [][]byte{
		testingutils.EncodeMessage(testingutils.PartialSignatureMsgForSlot(2, testingutils.DefaultMsgID, testingutils.TestingValidatorIndex, testingutils.DefaultKeySet)),
		testingutils.EncodeMessage(testingutils.PartialSignatureMsgForSlot(1, testingutils.DefaultMsgID, testingutils.TestingValidatorIndex, testingutils.DefaultKeySet)),
	}

	return &MessageValidationTest{
		Name:          "slot already advanced for partial signature message",
		Messages:      msgs,
		ExpectedError: validation.ErrSlotAlreadyAdvanced.Error(),
	}
}
