package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotProposer tests a justified round change but node is not the proposer
func NotProposer() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.Height = tests.ChangeProposerFuncInstanceHeight // will change proposer default for tests

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndHeight(ks.OperatorKeys[1], types.OperatorID(1), 2, tests.ChangeProposerFuncInstanceHeight),
		testingutils.TestingRoundChangeMessageWithRoundAndHeight(ks.OperatorKeys[2], types.OperatorID(2), 2, tests.ChangeProposerFuncInstanceHeight),
		testingutils.TestingRoundChangeMessageWithRoundAndHeight(ks.OperatorKeys[3], types.OperatorID(3), 2, tests.ChangeProposerFuncInstanceHeight),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "round change justification not proposer",
		Pre:           pre,
		PostRoot:      "cd38e73c2b2c13d4fc7a7ff60d596fbe274fd0ccaa1a19ad96b1f1520e95de1d",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingRoundChangeMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, tests.ChangeProposerFuncInstanceHeight,
				[32]byte{}, 0, [][]byte{}),
		},
	}
}
