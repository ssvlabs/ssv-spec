package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RoundChangeNotPreparedJustifications tests valid justified change round (not prev prepared)
func RoundChangeNotPreparedJustifications() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingRoundChangeMessageWithParams(
		ks.OperatorKeys[1], types.OperatorID(1), 10, qbft.FirstHeight, testingutils.TestingQBFTRootData, qbft.NoRound, nil)

	return &tests.MsgSpecTest{
		Name: "rc not prev prepared justifications",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
	}
}
