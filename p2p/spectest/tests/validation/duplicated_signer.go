package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicatedSigner tests a message with a duplicated signer
func DuplicatedSigner() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1, 2, 2}, // Duplicated signer
		Signatures:  [][]byte{testingutils.MockRSASignature[:], testingutils.MockRSASignature[:], testingutils.MockRSASignature[:]},
		SSVMessage:  &types.SSVMessage{},
	}

	return &MessageValidationTest{
		Name:          "duplicated signer",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrDuplicatedSigner.Error(),
	}
}
