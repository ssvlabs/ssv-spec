package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// SignedMessageSigner0 tests SignedMessage signer == 0
func SignedMessageSigner0() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.MultiSignAleaMsg(
		[]*bls.SecretKey{
			ks.Shares[1],
			ks.Shares[2],
			ks.Shares[3],
		},
		[]types.OperatorID{1, 2, 0},
		&alea.Message{
			MsgType:    alea.ProposalMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytesAlea([]byte{1, 2, 3, 4}),
		})

	return &tests.MsgSpecTest{
		Name: "signer 0",
		Messages: []*alea.SignedMessage{
			msg,
		},
		ExpectedError: "signer ID 0 not allowed",
	}
}
