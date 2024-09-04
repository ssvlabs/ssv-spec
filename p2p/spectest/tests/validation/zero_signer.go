package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ZeroSigner tests a message with siner zero
func ZeroSigner() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{0},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage:  &types.SSVMessage{},
	}

	return &MessageValidationTest{
		Name: "zero signer",

		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrZeroSigner.Error(),
	}
}
