package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// RoundChangeDataInvalidPreparedValue tests PreparedRound != NoRound && PreparedValue == nil
func RoundChangeDataInvalidPreparedValue() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingRoundChangeMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 10, qbft.FirstHeight, [32]byte{}, 2, nil)

	return &tests.MsgSpecTest{
		Name: "rc prepared no value",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "round change prepared value invalid",
	}
}
