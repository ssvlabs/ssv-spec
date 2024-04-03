package prepare

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MultiSigner tests prepare msg with > 1 signers
func MultiSigner() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.NetworkKeys[1], types.OperatorID(1))

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMultiSignerMessage(
			[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2]},
			[]types.OperatorID{1, 2},
		),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "prepare multi signer",
		Pre:           pre,
		PostRoot:      "b61f5233721865ca43afc68f4ad5045eeb123f6e8f095ce76ecf956dabc74713",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: msg allows 1 signer",
	}
}
