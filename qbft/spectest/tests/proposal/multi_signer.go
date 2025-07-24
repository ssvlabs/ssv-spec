package proposal

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MultiSigner tests a proposal msg with > 1 signers
func MultiSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingMultiSignerProposalMessage(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2]},
			[]types.OperatorID{1, 2},
		),
	}

	test := tests.NewMsgProcessingSpecTest(
		"proposal multi signer",
		testdoc.ProposalMultiSignerDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"invalid signed message: msg allows 1 signer",
		nil,
	)

	test.SetPrivateKeys(ks)

	return test
}
