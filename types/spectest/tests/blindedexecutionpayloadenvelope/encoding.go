package blindedexecutionpayloadenvelope

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BlindedExecutionPayloadEnvelopeEncoding tests encoding and decoding a BlindedExecutionPayloadEnvelope object.
func BlindedExecutionPayloadEnvelopeEncoding() *EncodingTest {
	byts, err := testingutils.TestBlindedExecutionPayloadEnvelope.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	root, err := testingutils.TestBlindedExecutionPayloadEnvelope.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"blinded execution payload envelope encoding",
		testdoc.BlindedExecutionPayloadEnvelopeEncodingTestDoc,
		byts,
		root,
	)
}
