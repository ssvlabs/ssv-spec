package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoPreviousProposal tests prepare msg without receiving a previous proposal state.ProposalAcceptedForCurrentRound == nil
func NoPreviousProposal() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "no previous proposal for prepare",
		Pre:           pre,
		PostRoot:      "620ad2417e47411537db8df9d4a072327e3c3efc391c3162867f30d5bf9af52c",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: did not receive proposal for this round",
	}
}
