package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignersNotSorted tests a message with non sorted signers
func SignersNotSorted() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{2, 1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:], testingutils.MockRSASignature[:]},
		SSVMessage:  &types.SSVMessage{},
	}

	return &MessageValidationTest{
		Name: "signers not sorted",

		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrSignersNotSorted.Error(),
	}
}
