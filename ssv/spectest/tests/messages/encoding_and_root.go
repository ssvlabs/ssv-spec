package messages

import (
	"github.com/bloxapp/ssv-spec/types"
)

// EncodingAndRoot tests SignedPartialSignatureMessage encoding + root
func EncodingAndRoot() *MsgSpecTest {
	msg := &types.SignedPartialSignatureMessage{
		Signature: make([]byte, 96),
		Signer:    12,
		Message: types.PartialSignatureMessages{
			Type: types.PostConsensusPartialSig,
			Messages: []*types.PartialSignatureMessage{
				{
					PartialSignature: make([]byte, 96),
					Signer:           12,
					SigningRoot:      make([]byte, 32),
				},
				{
					PartialSignature: make([]byte, 96),
					Signer:           12,
					SigningRoot:      make([]byte, 32),
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
		ExpectedRoots: [][]byte{
			r,
		},
	}
}
