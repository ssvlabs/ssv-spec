package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsg tests a duplicate prepare msg processing
func DuplicateMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))

	msgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "duplicate prepare message",
		Pre:           pre,
		PostRoot:      "9c77f5c787ae2ed292c7b1e837ca183cc4cad3316287b8173bbd043b5e764b06",
		InputMessages: msgs,
	}
}
