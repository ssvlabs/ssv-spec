package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// SignedMsgDuplicateSigners tests SignedMessage with duplicate signers
func SignedMsgDuplicateSigners() *tests.MsgSpecTest {
	msg := testingutils.MultiSignAleaMsg(
		[]*bls.SecretKey{
			testingutils.Testing4SharesSet().Shares[1],
			testingutils.Testing4SharesSet().Shares[1],
			testingutils.Testing4SharesSet().Shares[2],
		},
		[]types.OperatorID{1, 2, 3},
		&alea.Message{
			MsgType:    alea.ProposalMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytesAlea([]byte{1, 2, 3, 4}),
		})
	msg.Signers = []types.OperatorID{1, 1, 2}

	return &tests.MsgSpecTest{
		Name: "duplicate signers",
		Messages: []*alea.SignedMessage{
			msg,
		},
		ExpectedError: "non unique signer",
	}
}
