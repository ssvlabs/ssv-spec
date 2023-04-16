package messages

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
)

// EncodingAndRoot tests SignedPartialSignatureMessage encoding + root
func EncodingAndRoot() tests.SpecTest {
	msg := &types.SignedPartialSignatureMessage{
		Signature: make([]byte, 96),
		Signer:    12,
		Message: types.PartialSignatureMessages{
			Type: types.PostConsensusPartialSig,
			Messages: []*types.PartialSignatureMessage{
				{
					PartialSignature: make([]byte, 96),
					Signer:           12,
					SigningRoot:      [32]byte{},
				},
				{
					PartialSignature: make([]byte, 96),
					Signer:           12,
					SigningRoot:      [32]byte{},
				},
			},
		},
	}

	r, _ := msg.GetRoot()
	byts, _ := msg.Encode()

	return &MsgSpecTest{
		Name: "encoding",
		Messages: []*types.SignedPartialSignatureMessage{
			msg,
		},
		EncodedMessages: [][]byte{
			byts,
		},
		ExpectedRoots: [][32]byte{
			r,
		},
	}
}
