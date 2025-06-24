package beaconvote

import (
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BeaconVoteEncoding tests encoding and decoding a BeaconVote object
func BeaconVoteEncoding() *EncodingTest {

	bv := testingutils.TestBeaconVote

	byts, err := bv.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := bv.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"beacon vote encoding",
		byts,
		root,
	)
}
