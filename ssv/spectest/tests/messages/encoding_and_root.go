package messages

import (
	"github.com/bloxapp/ssv-spec/ssv"
)

// EncodingAndRoot tests SignedPartialSignatureMessage encoding + root
func EncodingAndRoot() *MsgSpecTest {
	msg := &ssv.SignedPartialSignatureMessage{
		Signature: make([]byte, 96),
		Signer:    12,
		Message: ssv.PartialSignatureMessages{
			Type: ssv.PostConsensusPartialSig,
			Messages: []*ssv.PartialSignatureMessage{
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
		Messages: []*ssv.SignedPartialSignatureMessage{
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
