package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SignedMsgSigner0 tests SignedPartialSignatureMessage signer == 0
func SignedMsgSigner0() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msgPre := testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)
	msgPre.Signer = 0
	msgPost := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msgPost.Signer = 0

	return &MsgSpecTest{
		Name: "signed message signer 0",
		Messages: []*types.SignedPartialSignatureMessage{
			msgPre,
			msgPost,
		},
		ExpectedError: "signer ID 0 not allowed",
	}
}
