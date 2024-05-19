package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
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
		PostRoot:      "6cd4ca943b96bbc598eb09ade7bcfc273759aa0b14a5b35199093383463a0458",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, tests.ChangeProposerFuncInstanceHeight,
				[32]byte{}, 0, [][]byte{}),
		},
	}
}
