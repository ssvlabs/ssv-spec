package signedbuilderrequestauth

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignedBuilderRequestAuthEncoding tests encoding and decoding a SignedBuilderRequestAuth object (SIP #94 §5).
func SignedBuilderRequestAuthEncoding() *EncodingTest {

	s := testingutils.TestSignedBuilderRequestAuth

	byts, err := s.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	root, err := s.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"signed builder request auth encoding",
		testdoc.SignedBuilderRequestAuthEncodingTestDoc,
		byts,
		root,
	)
}
