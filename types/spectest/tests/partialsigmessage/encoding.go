package partialsigmessage

import (
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Encoding tests encoding of a ssv message
func Encoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)

	byts, err := msg.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := msg.GetRoot()
	if err != nil {
		panic(err.Error())
	}

	return &EncodingTest{
		Name:         "SignedPartialSignatureMessage encoding",
		Data:         byts,
		ExpectedRoot: root,
	}
}
