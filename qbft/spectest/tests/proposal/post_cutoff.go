package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostCutoff tests processing a proposal msg when round >= cutoff
func PostCutoff() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.Round = 15

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 15),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round cutoff proposal message",
		Pre:           pre,
		PostRoot:      "256af2e6f30d3c1d3325161236c22ff0dde99a1985d74ad0f0a6f7da4aae888a",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
