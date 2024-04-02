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
		PostRoot:      "29e39b5903e17532811400d7ecec856640bd48d8a1dbae5cc118a234effea445",
		InputMessages: msgs,
	}
}
