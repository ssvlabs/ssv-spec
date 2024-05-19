package commit

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MultiSignerNoOverlap tests a multi signer commit msg which doesn't overlap previous valid commits
func MultiSignerNoOverlap() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], 1),

		testingutils.TestingPrepareMessage(ks.Shares[1], 1),
		testingutils.TestingPrepareMessage(ks.Shares[2], 2),
		testingutils.TestingPrepareMessage(ks.Shares[3], 3),

		testingutils.TestingCommitMessage(ks.Shares[1], 1),
		testingutils.TestingCommitMultiSignerMessage([]*bls.SecretKey{ks.Shares[2], ks.Shares[3]}, []types.OperatorID{2, 3}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "multi signer, no overlap",
		Pre:           pre,
		PostRoot:      "c72a82fa102b581273a5ba16dbadddee88fb850649fe7fb510986f126a07a04d",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], 1),
			testingutils.TestingCommitMessage(ks.Shares[1], 1),
		},
		ExpectedError: "invalid signed message: msg allows 1 signer",
	}
}
