package partialsigmessage

import (
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Encoding tests encoding of a ssv message
func Encoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, gloas.DataVersionGloas)

	byts, err := msg.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := msg.GetRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"PartialSignatureMessages encoding",
		testdoc.PartialSignatureMessageEncodingTestDoc,
		byts,
		root,
	)
}
