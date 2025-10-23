package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
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

	test := tests.NewMsgProcessingSpecTest(
		"wrong proposer",
		testdoc.ProposalWrongProposerDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		types.ProposalLeaderInvalidErrorCode,
		nil,
		ks,
	)

	return test
}
