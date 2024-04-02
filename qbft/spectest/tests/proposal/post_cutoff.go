package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostCutoff tests processing a proposal msg when round >= cutoff
func PostCutoff() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.Round = 15

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithRound(ks.Shares[1], types.OperatorID(1), 15),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round cutoff proposal message",
		Pre:           pre,
		PostRoot:      "95f83884e69cdfd6377a6b2b6af4d89d76902f07b76b445cf3fe5db27c8dc8b7",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
