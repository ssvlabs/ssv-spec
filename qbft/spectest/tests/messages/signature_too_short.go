package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// SignedMsgSigTooShort tests SignedMessage len(signature) < 96
func SignedMsgSigTooShort() *tests.MsgSpecTest {
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
	})
	msg.Signature = make([]byte, 95)

	return &tests.MsgSpecTest{
		Name: "signature too short",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "message signature is invalid",
	}
}
