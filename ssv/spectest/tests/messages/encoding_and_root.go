package messages

import (
	"github.com/bloxapp/ssv-spec/ssv"
)

// EncodingAndRoot tests SignedPartialSignatures encoding + root
func EncodingAndRoot() *MsgSpecTest {
	msg := &ssv.SignedPartialSignatures{
		Signature: make([]byte, 96),
		Signer:    11,
		PartialSignatures: ssv.PartialSignatures{
			{
				Slot:        11,
				Signature:   make([]byte, 96),
				SigningRoot: make([]byte, 32),
			},
			{
				Slot:        11,
				Signature:   make([]byte, 96),
				SigningRoot: make([]byte, 32),
			},
		},
	}

	r, _ := msg.GetRoot()
	byts, _ := msg.Encode()

	return &MsgSpecTest{
		Name: "encoding",
		Messages: []*ssv.SignedPartialSignatures{
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
