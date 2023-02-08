package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// MsgDataNil tests data == nil
func MsgDataNil() *tests.MsgSpecTest {
	msg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       nil,
	})

	return &tests.MsgSpecTest{
		Name: "msg data nil",
		Messages: []*alea.SignedMessage{
			msg,
		},
		ExpectedError: "message data is invalid",
	}
}
