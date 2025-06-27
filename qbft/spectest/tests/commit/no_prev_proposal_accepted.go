package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoPrevAcceptedProposal tests a commit msg received without a previous accepted proposal
func NoPrevAcceptedProposal() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = nil
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
	}

	return tests.NewMsgProcessingSpecTest(
		"no previous accepted proposal",
		"Test processing of a commit message without a previous accepted proposal",
		pre,
		"",
		nil,
		msgs,
		nil,
		"invalid signed message: did not receive proposal for this round",
		nil,
	)
}
