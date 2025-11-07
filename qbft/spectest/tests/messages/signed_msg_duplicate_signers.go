package messages

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignedMsgDuplicateSigners tests SignedMessage with duplicate signers
func SignedMsgDuplicateSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessage(
		[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[1], ks.OperatorKeys[2]},
		[]types.OperatorID{1, 2, 3},
	)
	msg.OperatorIDs = []types.OperatorID{1, 1, 2}

	test := tests.NewMsgSpecTest(
		"duplicate signers",
		testdoc.MessagesSignedMsgDuplicateSignersDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		types.NonUniqueSignerErrorCode,
		ks,
	)

	return test
}
