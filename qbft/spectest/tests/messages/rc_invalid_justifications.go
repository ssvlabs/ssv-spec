package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// RoundChangeDataInvalidJustifications tests PreparedRound != NoRound len(RoundChangeJustification) == 0
func RoundChangeDataInvalidJustifications() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingRoundChangeMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 10, qbft.FirstHeight, testingutils.TestingQBFTRootData, 1, nil)

	return &tests.MsgSpecTest{
		Name: "rc prev prepared no justifications",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "round change justification invalid",
	}
}
