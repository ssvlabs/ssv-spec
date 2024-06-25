package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// SignedMsgMultiSigners tests SignedMessage with multi signers
func SignedMsgMultiSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingCommitMultiSignerMessage(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 3},
	)

	return &tests.MsgSpecTest{
		Name: "multi signers",
		Messages: []*qbft.SignedMessage{
			msg,
		},
	}
}
