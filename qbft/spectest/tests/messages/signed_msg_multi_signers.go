package messages

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SignedMsgMultiSigners tests SignedMessage with multi signers
func SignedMsgMultiSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingCommitMultiSignerMessage(
		[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
		[]types.OperatorID{1, 2, 3},
	)

	return &tests.MsgSpecTest{
		Name: "multi signers",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
	}
}
