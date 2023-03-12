package share

import "github.com/bloxapp/ssv-spec/types/testingutils"

// Encoding tests encoding of a ssv message
func Encoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)
	share.Graffiti = []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

	byts, err := share.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := share.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &EncodingTest{
		Name:         "share encoding",
		Data:         byts,
		ExpectedRoot: root,
	}
}
