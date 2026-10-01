package builderrequestauth

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthEncoding tests encoding and decoding a BuilderRequestAuth object (SIP #94 §5).
func BuilderRequestAuthEncoding() *EncodingTest {

	b := testingutils.TestBuilderRequestAuth

	byts, err := b.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	root, err := b.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"builder request auth encoding",
		testdoc.BuilderRequestAuthEncodingTestDoc,
		byts,
		root,
	)
}
