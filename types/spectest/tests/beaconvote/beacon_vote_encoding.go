package beaconvote

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
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
		testdoc.BeaconVoteEncodingTestDoc,
		byts,
		root,
	)
}

// BeaconVotePreGloasEncoding tests that the pre-Gloas encoding of a vote does not decode (SIP #94 §2): it is the
// same vote without its AttestationDataIndex, 112 bytes to the Gloas vote's 120, and the two encodings reject
// each other on length.
func BeaconVotePreGloasEncoding() *EncodingTest {
	test := NewEncodingTest(
		"beacon vote pre-gloas encoding",
		testdoc.BeaconVotePreGloasEncodingTestDoc,
		testingutils.TestPreGloasBeaconVoteByts,
		[32]byte{},
	)
	test.DecodeFails = true
	return test
}
