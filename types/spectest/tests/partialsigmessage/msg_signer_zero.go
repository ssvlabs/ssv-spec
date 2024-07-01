package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MessageSigner0 tests PartialSignatureMessage signer == 0
func MessageSigner0() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msgPre := testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)
	msgPre.Messages[0].Signer = 0
	msgPost := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msgPost.Messages[0].Signer = 0

	return &MsgSpecTest{
		Name: "message signer 0",
		Messages: []*types.PartialSignatureMessages{
			msgPre,
			msgPost,
		},
		ExpectedError: "message invalid: signer ID 0 not allowed",
	}
}
