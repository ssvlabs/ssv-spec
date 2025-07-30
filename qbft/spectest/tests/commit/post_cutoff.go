package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostCutoff tests processing a commit msg when round >= cutoff
func PostCutoff() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.Round = 15

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 15),
	}

	test := tests.NewMsgProcessingSpecTest(
		"round cutoff commit message",
		testdoc.CommitTestPostCutoffDoc,
		pre,
		"",
		nil,
		msgs,
		nil,
		"instance stopped processing messages",
		nil,
		ks,
	)

	return test
}
