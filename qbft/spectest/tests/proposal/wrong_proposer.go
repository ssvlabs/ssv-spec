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
		PostRoot:      "01489f7af13579b66ce3da156d4d10208c85a10365380f04e7b8d82d0a9679ce",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: proposal leader invalid",
	}
}
