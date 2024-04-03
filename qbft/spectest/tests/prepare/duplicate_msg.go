package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsg tests a duplicate prepare msg processing
func DuplicateMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.NetworkKeys[1], types.OperatorID(1))

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "duplicate prepare message",
		Pre:           pre,
		PostRoot:      "2e51f4773ca9d8091d39dca8bbc3bc4989404992cdba8d4044bfd952c2b8335a",
		InputMessages: msgs,
	}
}
