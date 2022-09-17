package messages

import (
	"github.com/bloxapp/ssv-spec/ssv"
)

// EncodingAndRoot tests SignedPartialSignature encoding + root
func EncodingAndRoot() *MsgSpecTest {
	msg := &ssv.SignedPartialSignature{
		Signature: make([]byte, 96),
		Signer:    11,
		Message: ssv.PartialSignatures{
			Messages: []*ssv.PartialSignature{
				{
					Slot:             11,
					PartialSignature: make([]byte, 96),
					Signer:           12,
					SigningRoot:      make([]byte, 32),
				},
				{
					Slot:             11,
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
		Messages: []*ssv.SignedPartialSignature{
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
