package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongHeight tests a proposal msg received with the wrong height
func WrongHeight() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithHeight(ks.Shares[1], types.OperatorID(1), 2),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "wrong proposal height",
		Pre:           pre,
		PostRoot:      "5b18ca0b470208d8d247543306850618f02bddcbaa7c37eb6d5b36eb3accb5fb",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg height",
	}
}
