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
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[2], types.OperatorID(2)),
	}

	return tests.NewMsgProcessingSpecTest(
		"wrong proposer",
		"Test proposal by a node that is not the designated proposer for the current round, expecting validation error.",
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"invalid signed message: proposal leader invalid",
		nil,
	)
}
