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
		PostRoot:      "367ef26442cb2fc12f7764227ee5bcf5af6155a78f5e99475b012732257403da",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingRoundChangeMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, tests.ChangeProposerFuncInstanceHeight,
				[32]byte{}, 0, [][]byte{}),
		},
	}
}
