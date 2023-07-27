package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ForceStop tests processing a proposal msg when instance force stopped
func ForceStop() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.ForceStop()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "force stop proposal message",
		Pre:           pre,
		PostRoot:      "5b18ca0b470208d8d247543306850618f02bddcbaa7c37eb6d5b36eb3accb5fb",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
