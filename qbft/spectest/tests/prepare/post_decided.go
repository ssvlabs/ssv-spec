package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests processing prepare msg after instance decided
func PostDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], 2),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], 3),

		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingCommitMessage(ks.OperatorKeys[2], 2),
		testingutils.TestingCommitMessage(ks.OperatorKeys[3], 3),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[4], 4),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare post decided",
		Pre:           pre,
		PostRoot:      "1fcd62b476d2e41e3541fdb41077dbcf707b977ee9134d64eaa19e5fb4f66368",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
			testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		},
	}
}
