package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnknownSigner tests a single commit received with an unknown signer
func UnknownSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], 2),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], 3),

		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 5),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "unknown commit signer",
		Pre:           pre,
		InputMessages: msgs,
		ExpectedError: "invalid signed message: signer not in committee",
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
			testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		},
	}
}
