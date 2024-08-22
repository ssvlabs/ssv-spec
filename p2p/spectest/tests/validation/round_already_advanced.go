package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RoundAlreadyAdvanced tests two consensus message with different rounds but the second message has a lower round value
func RoundAlreadyAdvanced() tests.SpecTest {
	return &MessageValidationTest{
		Name: "round already advanced",
		Messages: [][]byte{
			testingutils.EncodeMessage(testingutils.ConsensusMessageForRound(2, testingutils.DefaultMsgID)),
			testingutils.EncodeMessage(testingutils.ConsensusMessageForRound(1, testingutils.DefaultMsgID)),
		},
		ExpectedError: validation.ErrRoundAlreadyAdvanced.Error(),
	}
}
