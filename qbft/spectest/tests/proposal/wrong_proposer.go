package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongProposer tests a proposal by the wrong proposer
func WrongProposer() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[2], types.OperatorID(2)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "wrong proposer",
		Pre:           pre,
		InputMessages: msgs,
		ExpectedError: "invalid signed message: proposal leader invalid",
	}
}
