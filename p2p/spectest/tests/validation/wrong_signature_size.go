package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongRSASignatureSize tests a message with a signature with wrong size
func WrongRSASignatureSize() tests.SpecTest {

	signatureWithWrongSize := [4]byte{1, 2, 3, 4}

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{signatureWithWrongSize[:]},
		SSVMessage:  &types.SSVMessage{},
	}

	return &MessageValidationTest{
		Name: "wrong signature size",

		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrWrongRSASignatureSize.Error(),
	}
}
