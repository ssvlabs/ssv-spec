package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongHeight tests a proposal msg received with the wrong height
func WrongHeight() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithHeight(ks.OperatorKeys[1], types.OperatorID(1), 2),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "wrong proposal height",
		Pre:           pre,
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg height",
	}
}
