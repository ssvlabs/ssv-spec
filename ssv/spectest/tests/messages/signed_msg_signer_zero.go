package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// SignedMsgSigner0 tests SignedPartialSignatureMessage signer == 0
func SignedMsgSigner0() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msgPre := testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)
	msgPre.Signer = 0
	msgPost := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msgPost.Signer = 0

	return &MsgSpecTest{
		Name: "signed message signer 0",
		Messages: []*ssv.SignedPartialSignatureMessage{
			msgPre,
			msgPost,
		},
		ExpectedError: "signer ID 0 not allowed",
	}
}
