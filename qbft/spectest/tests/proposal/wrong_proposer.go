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
		PostRoot:      "613745b592755d889d7fdec2b3a7e3b54ff8b5d981bf1a81683f3804f3350727",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: proposal leader invalid",
	}
}
