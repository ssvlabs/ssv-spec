package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SignedMsgSigTooShort tests SignedMessage len(signature) < 96
func SignedMsgSigTooShort() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1))
	msg.Signature = make([]byte, 95)

	return &tests.MsgSpecTest{
		Name: "signature too short",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "message signature is invalid",
	}
}
