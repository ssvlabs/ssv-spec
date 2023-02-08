package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// SignedMsgNoSigners tests SignedMessage len(signers) == 0
func SignedMsgNoSigners() *tests.MsgSpecTest {
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
	})
	msg.Signers = nil

	return &tests.MsgSpecTest{
		Name: "no signers",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "message signers is empty",
	}
}
