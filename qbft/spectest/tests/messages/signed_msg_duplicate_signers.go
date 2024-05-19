package messages

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
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
