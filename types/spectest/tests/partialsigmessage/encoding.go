package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Encoding tests encoding of a ssv message
func Encoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionPhase0)

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
