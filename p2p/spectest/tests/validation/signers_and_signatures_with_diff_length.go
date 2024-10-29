package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignersAndSignaturesWithDifferentLength tests a message with len(signers) != len(signatures)
func SignersAndSignaturesWithDifferentLength() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1, 2},
		Signatures:  [][]byte{testingutils.MockRSASignature[:], testingutils.MockRSASignature[:], testingutils.MockRSASignature[:]},
		SSVMessage:  &types.SSVMessage{},
	}

	return &MessageValidationTest{
		Name: "signers and signatures with different length",

		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrSignersAndSignaturesWithDifferentLength.Error(),
	}
}
