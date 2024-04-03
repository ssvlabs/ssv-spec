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
		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "no previous proposal for prepare",
		Pre:           pre,
		PostRoot:      "5b18ca0b470208d8d247543306850618f02bddcbaa7c37eb6d5b36eb3accb5fb",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: did not receive proposal for this round",
	}
}
