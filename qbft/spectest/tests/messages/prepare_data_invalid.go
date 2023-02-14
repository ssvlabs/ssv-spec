package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PrepareDataInvalid tests prepare data is invalid
func PrepareDataInvalid() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingPrepareMessageWithParams(ks.Shares[1], types.OperatorID(1), qbft.FirstRound, qbft.FirstHeight, testingutils.WrongRoot)

	return &tests.MsgSpecTest{
		Name: "prepare data invalid",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "PrepareData data is invalid",
	}
}
