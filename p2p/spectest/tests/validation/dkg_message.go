package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DKGMessage tests message with the DKG msg type
func DKGMessage() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.DKGMsgType, // DKG msg type
			MsgID:   testingutils.DefaultMsgID,
			Data:    testingutils.EncodeQbftMessage(&qbft.Message{}),
		},
	}

	return &MessageValidationTest{
		Name:          "dkg message",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrDKGMessage.Error(),
	}
}
