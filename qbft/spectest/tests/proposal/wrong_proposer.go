package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongProposer tests a proposal by the wrong proposer
func WrongProposer() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[2], types.OperatorID(2)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "wrong proposer",
		Pre:           pre,
		InputMessages: msgs,
		ExpectedError: "invalid signed message: proposal leader invalid",
	}
}
