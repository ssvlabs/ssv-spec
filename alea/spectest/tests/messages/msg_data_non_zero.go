package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// MsgDataNonZero tests len(data) == 0
func MsgDataNonZero() *tests.MsgSpecTest {
	msg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       []byte{},
	})

	return &tests.MsgSpecTest{
		Name: "msg data len 0",
		Messages: []*alea.SignedMessage{
			msg,
		},
		ExpectedError: "message data is invalid",
	}
}
