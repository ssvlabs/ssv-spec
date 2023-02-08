package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// MessageSigner0 tests PartialSignatureMessage signer == 0
func MessageSigner0() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msgPre := testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)
	msgPre.Message.Messages[0].Signer = 0
	msgPre.Signer = 0
	msgPost := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msgPost.Message.Messages[0].Signer = 0
	msgPost.Signer = 0

	return &MsgSpecTest{
		Name: "message signer 0",
		Messages: []*ssv.SignedPartialSignatureMessage{
			msgPre,
			msgPost,
		},
		ExpectedError: "signer ID 0 not allowed",
	}
}
