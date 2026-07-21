package payloadattestationmessage

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PayloadAttestationMessageEncoding tests encoding and decoding a PayloadAttestationMessage object (SIP #94 §3).
func PayloadAttestationMessageEncoding() *EncodingTest {

	m := testingutils.TestPayloadAttestationMessage

	byts, err := m.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	root, err := m.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"payload attestation message encoding",
		testdoc.PayloadAttestationMessageEncodingTestDoc,
		byts,
		root,
	)
}
