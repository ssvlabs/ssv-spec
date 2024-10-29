package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UndecodableRoundChangeJustification tests a message with an undecodable round change justification
func UndecodableRoundChangeJustification() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   testingutils.TestingMessageID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				RoundChangeJustification: [][]byte{{1, 2, 3, 4}},
			}),
		},
	}

	return &MessageValidationTest{
		Name: "undecodable round change justification",

		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrUndecodableData.Error(),
	}
}
