package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SecondProposalForRound tests a second proposal (by same signer) for current round. state.ProposalAcceptedForCurrentRound != nil && signedProposal.Message.Round == state.Round
func SecondProposalForRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		// TODO: originally using different value
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "second proposal for round",
		Pre:           pre,
		PostRoot:      "a218ef33486e9bdebef968f145f69a3e4ec9c4d824842297d3f6510d72ed5003",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		},
		ExpectedError: "invalid signed message: proposal is not valid with current state",
	}
}
