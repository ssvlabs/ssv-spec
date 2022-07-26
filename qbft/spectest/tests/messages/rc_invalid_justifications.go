package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// RoundChangeDataInvalidJustifications tests PreparedRound != NoRound len(RoundChangeJustification) == 0
func RoundChangeDataInvalidJustifications() *tests.MsgSpecTest {
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.RoundChangeMsgType,
		Height:     qbft.FirstHeight,
		Round:      10,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, 1, nil),
	})

	return &tests.MsgSpecTest{
		Name: "rc prev prepared no justifications",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "round change justification invalid",
	}
}
