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
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "duplicate prepare message",
		Pre:           pre,
		PostRoot:      "be7d7c5c8880746e64eaa79f164898803b0832b8d56e19e38e1e8edbe985ae74",
		InputMessages: msgs,
	}
}
