package commit

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MultiSignerNoOverlap tests a multi signer commit msg which doesn't overlap previous valid commits
func MultiSignerNoOverlap() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.NetworkKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], 2),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[3], 3),

		testingutils.TestingCommitMessage(ks.NetworkKeys[1], 1),
		testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.NetworkKeys[2], ks.NetworkKeys[3]}, []types.OperatorID{2, 3}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "multi signer, no overlap",
		Pre:           pre,
		PostRoot:      "4e569d9a6c0421d2bb69a4c544f8f1e67c73a129d4e6bd1304ddbae8812cfa38",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.NetworkKeys[1], 1),
			testingutils.TestingCommitMessage(ks.NetworkKeys[1], 1),
		},
		ExpectedError: "invalid signed message: msg allows 1 signer",
	}
}
