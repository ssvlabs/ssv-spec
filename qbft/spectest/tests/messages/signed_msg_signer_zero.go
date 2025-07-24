package messages

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignedMessageSigner0 tests SignedMessage signer == 0
func SignedMessageSigner0() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessage(
		[]*rsa.PrivateKey{
			ks.OperatorKeys[1],
			ks.OperatorKeys[2],
			ks.OperatorKeys[3],
		},
		[]types.OperatorID{1, 2, 0},
	)

	test := tests.NewMsgSpecTest(
		"signer 0",
		testdoc.MessagesSignedMsgSignerZeroDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		"signer ID 0 not allowed",
	)

	test.SetPrivateKeys(ks)

	return test
}
