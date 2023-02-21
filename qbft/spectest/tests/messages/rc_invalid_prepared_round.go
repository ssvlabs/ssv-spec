package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// RoundChangeDataInvalidPreparedRound tests PreparedValue != nil && PreparedRound == NoRound
func RoundChangeDataInvalidPreparedRound() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 10)

	return &tests.MsgSpecTest{
		Name: "rc prev prepared no round",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "round change justification invalid",
	}
}
