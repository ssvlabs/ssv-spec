package messages

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignedMessageEmptySignature tests an invalid SignedSSVMessage with an empty signature
func SignedMessageEmptySignature() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessage(
		[]*rsa.PrivateKey{
			ks.OperatorKeys[1],
			ks.OperatorKeys[2],
			ks.OperatorKeys[3],
		},
		[]types.OperatorID{1, 2, 3},
	)

	msg.Signatures[0] = make([]byte, 0)

	return &tests.MsgSpecTest{
		Name: "signedssvmessage with empty signature",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		ExpectedError: "empty signature",
	}
}
