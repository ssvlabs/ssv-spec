package commit

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests processing a commit msg after instance decided
func PostDecided() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1)

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingCommitMessage(ks.OperatorKeys[2], 2),
		testingutils.TestingCommitMessage(ks.OperatorKeys[3], 3),
		testingutils.TestingCommitMessage(ks.OperatorKeys[4], 4),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "post decided",
		Pre:            pre,
		PostRoot:       "eed0ce97278b3b003f03c3b9eadea62611a413638be9fec790188c5c7f6e67a7",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
	}
}
