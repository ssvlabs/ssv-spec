package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignedMsgNoSigners tests SignedMessage len(signers) == 0
func SignedMsgNoSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1))
	msg.OperatorIDs = nil

	test := tests.NewMsgSpecTest(
		"no signers",
		testdoc.MessagesSignedMsgNoSignersDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		types.NoSignersErrorCode,
		ks,
	)

	return test
}
