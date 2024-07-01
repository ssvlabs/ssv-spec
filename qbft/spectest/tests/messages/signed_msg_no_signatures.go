package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignedMessageNoSignatures tests an invalid SignedSSVMessage with no signatures
func SignedMessageNoSignatures() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1)
	msg.Signatures = make([][]byte, 0)

	return &tests.MsgSpecTest{
		Name: "signedssvmessage with no signatures",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		ExpectedError: "no signatures",
	}
}
