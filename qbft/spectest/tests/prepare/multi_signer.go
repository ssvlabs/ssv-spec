package prepare

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MultiSigner tests prepare msg with > 1 signers
func MultiSigner() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMultiSignerMessage(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2]},
			[]types.OperatorID{1, 2},
		),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "prepare multi signer",
		Pre:           pre,
		PostRoot:      "d3c3540b20d61319771ef7b01c26d172fbabeabad5c09b1415c72a34e0fd145f",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: msg allows 1 signer",
	}
}
