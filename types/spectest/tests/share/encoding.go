package share

import "github.com/ssvlabs/ssv-spec/types/testingutils"

// Encoding tests encoding of a ssv message
func Encoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks, testingutils.TestingValidatorIndex)

	byts, err := share.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := share.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"share encoding",
		byts,
		root,
	)
}
