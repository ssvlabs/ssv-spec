package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotProposer tests a justified round change but node is not the proposer
func NotProposer() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.Height = tests.ChangeProposerFuncInstanceHeight // will change proposer default for tests

	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndHeight(ks.Shares[1], types.OperatorID(1), 2, tests.ChangeProposerFuncInstanceHeight),
		testingutils.TestingRoundChangeMessageWithRoundAndHeight(ks.Shares[2], types.OperatorID(2), 2, tests.ChangeProposerFuncInstanceHeight),
		testingutils.TestingRoundChangeMessageWithRoundAndHeight(ks.Shares[3], types.OperatorID(3), 2, tests.ChangeProposerFuncInstanceHeight),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "round change justification not proposer",
		Pre:           pre,
		PostRoot:      "67e2bb70161179228a25db083efcba16569d08d588d86ace24a96522e5bd01d3",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, tests.ChangeProposerFuncInstanceHeight,
				[32]byte{}, 0, [][]byte{}),
		},
	}
}
