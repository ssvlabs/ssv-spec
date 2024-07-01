package messages

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignedMessageNilSSVMessage tests an invalid SignedSSVMessage with a nil SSVMessage
func SignedMessageNilSSVMessage() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessage(
		[]*rsa.PrivateKey{
			ks.OperatorKeys[1],
			ks.OperatorKeys[2],
			ks.OperatorKeys[3],
		},
		[]types.OperatorID{1, 2, 3},
	)

	msg.SSVMessage = nil

	return &tests.MsgSpecTest{
		Name: "signedssvmessage with nil ssvmessage",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		ExpectedError: "nil SSVMessage",
	}
}
