package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoSigners tests a message with no signers
func NoSigners() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: make([]types.OperatorID, 0), // No signers
		Signatures:  [][]byte{{1, 2, 3, 4}},
		SSVMessage:  &types.SSVMessage{},
	}

	return &MessageValidationTest{
		Name:          "no signers",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrNoSigners.Error(),
	}
}
