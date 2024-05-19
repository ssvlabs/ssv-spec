package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicateMsg tests a duplicate commit msg processing
func DuplicateMsg() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], 1)

	msgs := []*qbft.SignedMessage{
		testingutils.TestingCommitMessage(ks.Shares[1], 1),
		testingutils.TestingCommitMessage(ks.Shares[1], 1),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "duplicate commit message",
		Pre:           pre,
		PostRoot:      "865b610062195fe9635487cfed19d3fc40ebffa96e3d440ccdf1209e9b93186d",
		InputMessages: msgs,
	}
}
