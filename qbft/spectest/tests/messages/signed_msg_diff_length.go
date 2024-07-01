package messages

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignedMessageDifferentLength tests an invalid SignedSSVMessage with different number of signers and signatures
func SignedMessageDifferentLength() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessage(
		[]*rsa.PrivateKey{
			ks.OperatorKeys[1],
			ks.OperatorKeys[2],
			ks.OperatorKeys[3],
		},
		[]types.OperatorID{1, 2, 3},
	)

	msg.OperatorIDs = []types.OperatorID{1, 2, 3, 4}

	return &tests.MsgSpecTest{
		Name: "signedssvmessage with different length of signers and signatures",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		ExpectedError: "number of signatures is different than number of signers",
	}
}
