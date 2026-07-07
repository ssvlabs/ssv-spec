package gloasbeaconvote

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasBeaconVoteEncoding tests encoding and decoding a GloasBeaconVote object (SIP #94 §2).
func GloasBeaconVoteEncoding() *EncodingTest {

	bv := testingutils.TestGloasBeaconVote

	byts, err := bv.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := bv.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"gloas beacon vote encoding",
		testdoc.GloasBeaconVoteEncodingTestDoc,
		byts,
		root,
	)
}
