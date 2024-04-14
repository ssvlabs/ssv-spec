package commit

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsg tests a duplicate commit msg processing
func DuplicateMsg() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1)

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "duplicate commit message",
		Pre:           pre,
		PostRoot:      "d5835a13c15b8433c235ebda918a844660e4f17b7873b7ed738c16bc8e6004e0",
		InputMessages: msgs,
	}
}
