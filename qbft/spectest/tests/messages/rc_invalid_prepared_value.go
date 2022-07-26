package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// RoundChangeDataInvalidPreparedValue tests PreparedRound != NoRound && PreparedValue == nil
func RoundChangeDataInvalidPreparedValue() *tests.MsgSpecTest {
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.RoundChangeMsgType,
		Height:     qbft.FirstHeight,
		Round:      10,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.RoundChangePreparedDataBytes(nil, 2, nil),
	})

	return &tests.MsgSpecTest{
		Name: "rc prepared no value",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "round change prepared value invalid",
	}
}
