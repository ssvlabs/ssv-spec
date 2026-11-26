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

// PayloadAttestationDataInvalidBoolean tests that a boolean encoded as a byte other than 0 or 1 does not decode
// (SSZ allows no other value): the fixture data with its payload_present byte set to 2.
func PayloadAttestationDataInvalidBoolean() *EncodingTest {
	byts, err := testingutils.TestPayloadAttestationData.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	// beacon_block_root (32 bytes) and slot (8 bytes) come before payload_present.
	const payloadPresentOffset = 40
	byts[payloadPresentOffset] = 2

	test := NewEncodingTest(
		"payload attestation data invalid boolean",
		testdoc.PayloadAttestationDataInvalidBooleanTestDoc,
		byts,
		[32]byte{},
	)
	test.DecodeFails = true
	return test
}
