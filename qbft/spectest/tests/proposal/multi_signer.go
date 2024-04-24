package proposal

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
		PostRoot:       "620ad2417e47411537db8df9d4a072327e3c3efc391c3162867f30d5bf9af52c",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: msg allows 1 signer",
	}
}
