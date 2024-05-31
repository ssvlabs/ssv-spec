package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
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
		PostRoot:      "af5b177d7c1588c741a550bd8100f0099a8a85e9249589d15750192bccf4e631",
		InputMessages: msgs,
	}
}
