package proposal

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MultiSigner tests a proposal msg with > 1 signers
func MultiSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingMultiSignerProposalMessage(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2]},
			[]types.OperatorID{1, 2},
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "proposal multi signer",
		Pre:            pre,
		PostRoot:       "3d11aa7331a7aa79d3403ac1af61569f1eae0547f54f15dca7e9e07b1ab0573d",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: msg allows 1 signer",
	}
}
