package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
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
		PostRoot:      "fc38640f9cf1613bb500f2a9b8aacfd2685727adc5a261d28646f494f505b357",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
