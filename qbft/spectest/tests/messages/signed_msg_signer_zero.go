package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// SignedMessageSigner0 tests SignedMessage signer == 0
func SignedMessageSigner0() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{
			ks.Shares[1],
			ks.Shares[2],
			ks.Shares[3],
		},
		[]types.OperatorID{1, 2, 0},
		&qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		})

	return &tests.MsgSpecTest{
		Name: "signer 0",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "signer ID 0 not allowed",
	}
}
