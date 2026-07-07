package envelopeconsensusdata

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// EnvelopeConsensusDataEncoding tests encoding and decoding an EnvelopeConsensusData object (SIP #94 §6).
func EnvelopeConsensusDataEncoding() *EncodingTest {

	cd := testingutils.TestEnvelopeConsensusData

	byts, err := cd.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := cd.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"envelope consensus data encoding",
		testdoc.EnvelopeConsensusDataEncodingTestDoc,
		byts,
		root,
	)
}
