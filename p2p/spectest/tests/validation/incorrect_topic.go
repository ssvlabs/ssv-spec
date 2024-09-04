package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// IncorrectTopic tests message in the wrong topic
func IncorrectTopic() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   testingutils.DefaultMsgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				Round:      qbft.FirstRound,
				Identifier: testingutils.DefaultMsgID[:],
			}),
		},
	}

	return &MessageValidationTest{
		Name:          testingutils.TestingWrongTopic, // Wrong topic
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrIncorrectTopic.Error(),
		Topic:         testingutils.TestingWrongTopic,
	}
}
