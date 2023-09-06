package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongSig tests a signed round change msg with wrong signature
func WrongSig() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(2), 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change invalid sig",
		Pre:            pre,
		PostRoot:       "3ad53d1ce5f9ccbcf3e2a37402ee5f08f9e8113042c6ed5b42045dc9dcc1844a",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: msg signature invalid: failed to verify signature",
	}
}
