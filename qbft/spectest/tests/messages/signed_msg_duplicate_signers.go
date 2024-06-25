package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// SignedMsgDuplicateSigners tests SignedMessage with duplicate signers
func SignedMsgDuplicateSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessage(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[1], ks.Shares[2]},
		[]types.OperatorID{1, 2, 3},
	)
	msg.Signers = []types.OperatorID{1, 1, 2}

	return &tests.MsgSpecTest{
		Name: "duplicate signers",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "non unique signer",
	}
}
