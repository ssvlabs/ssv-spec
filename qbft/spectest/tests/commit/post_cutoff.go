package commit

import (
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

	return &tests.MsgProcessingSpecTest{
		Name:          "round cutoff commit message",
		Pre:           pre,
		PostRoot:      "d9306681b75a11b98919b642b929a95db1c6b3de653e27a9cf840b32fcdeb9fe",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
