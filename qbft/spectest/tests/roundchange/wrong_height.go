package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongHeight tests a round change msg with wrong height
func WrongHeight() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndHeight(ks.OperatorKeys[1], types.OperatorID(1), 2, 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change invalid height",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: wrong msg height",
	}
}
