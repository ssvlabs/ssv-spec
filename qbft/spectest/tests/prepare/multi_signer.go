package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiSigner tests prepare msg with > 1 signers
func MultiSigner() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))

	msgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMultiSignerMessage(
			[]*bls.SecretKey{ks.Shares[1], ks.Shares[2]},
			[]types.OperatorID{1, 2},
		),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "prepare multi signer",
		Pre:           pre,
		PostRoot:      "35f7aaa4f57445a48f189b8c8edf66a3d0e3d54fde910097b6946ca6fa4d73ab",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: msg allows 1 signer",
	}
}
