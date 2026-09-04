package envelopedissemination

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// EnvelopeDisseminationEncoding tests encoding and decoding an EnvelopeDissemination object (SIP #94 §6).
func EnvelopeDisseminationEncoding() *EncodingTest {

	d := testingutils.TestEnvelopeDissemination

	byts, err := d.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := d.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"envelope dissemination encoding",
		testdoc.EnvelopeDisseminationEncodingTestDoc,
		byts,
		root,
	)
}
