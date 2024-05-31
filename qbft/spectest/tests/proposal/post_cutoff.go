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
		PostRoot:      "962a657645159c8371a9aed54e55637cbde19d88db7bf05126c96e90653cb140",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
