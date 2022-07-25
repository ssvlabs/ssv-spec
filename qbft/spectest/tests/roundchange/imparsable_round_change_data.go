package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ImparsableRoundChangeData tests a round change msg data imparsable to round change data struct
func ImparsableRoundChangeData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       []byte{1, 1},
		}),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change data imparsable",
		Pre:            pre,
		PostRoot:       "4aafcc4aa9e2435579c85aa26e659fe650aefb8becb5738d32dd9286f7ff27c3",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "round change msg invalid: could not get roundChange data : could not decode round change data from message: invalid character '\\x01' looking for beginning of value",
	}
}
