package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialSignatureSSVMessageUndecodableData tests a partial signature message with a SSVMessage that has an undecodable data
func PartialSignatureSSVMessageUndecodableData() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   testingutils.TestingMessageID,
			Data:    []byte{1, 2, 3, 4}, // Undecodable data
		},
	}

	return &MessageValidationTest{
		Name:          "partial signature ssvmesage undecodable data",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrUndecodableData.Error(),
	}
}
