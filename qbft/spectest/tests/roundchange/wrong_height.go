package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
		PostRoot:       "39d85fcbfb9564563798e3e4990f8f55efa6c86bd22d7d0ad0dfd3643e585954",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: wrong msg height",
	}
}
