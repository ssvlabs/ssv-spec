package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostCutoff tests processing a round change msg when round >= cutoff
func PostCutoff() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.Round = 15

	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 16),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round cutoff round change message",
		Pre:           pre,
		PostRoot:      "5efc161aa6053d5137d6451343e1045708d486a3a9ecb0c08baaad76da57da16",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
