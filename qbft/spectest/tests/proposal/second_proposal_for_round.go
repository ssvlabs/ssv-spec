package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
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
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		},
		ExpectedError: "invalid signed message: proposal is not valid with current state",
	}
}
