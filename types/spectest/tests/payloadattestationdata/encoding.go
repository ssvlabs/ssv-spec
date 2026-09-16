package payloadattestationdata

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PayloadAttestationDataEncoding tests encoding and decoding a PayloadAttestationData object (SIP #94 §3).
func PayloadAttestationDataEncoding() *EncodingTest {

	d := testingutils.TestPayloadAttestationData

	byts, err := d.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	root, err := d.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"payload attestation data encoding",
		testdoc.PayloadAttestationDataEncodingTestDoc,
		byts,
		root,
	)
}
