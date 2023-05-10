package messages

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MessageSigner0 tests PartialSignatureMessage signer == 0
func MessageSigner0() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msgPre := testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix)
	msgPre.Message.Messages[0].Signer = 0
	msgPre.Signer = 0
	msgPost := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msgPost.Message.Messages[0].Signer = 0
	msgPost.Signer = 0

	return &MsgSpecTest{
		Name: "message signer 0",
		Messages: []*types.SignedPartialSignatureMessage{
			msgPre,
			msgPost,
		},
		ExpectedError: "signer ID 0 not allowed",
	}
}
