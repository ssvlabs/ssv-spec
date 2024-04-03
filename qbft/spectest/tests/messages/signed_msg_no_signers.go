package messages

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SignedMsgNoSigners tests SignedMessage len(signers) == 0
func SignedMsgNoSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingCommitMessage(ks.NetworkKeys[1], types.OperatorID(1))
	msg.OperatorID = nil

	return &tests.MsgSpecTest{
		Name: "no signers",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		ExpectedError: "message signers is empty",
	}
}
