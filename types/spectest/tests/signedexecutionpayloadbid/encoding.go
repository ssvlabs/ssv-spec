package signedexecutionpayloadbid

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignedExecutionPayloadBidEncoding tests encoding and decoding a SignedExecutionPayloadBid object.
func SignedExecutionPayloadBidEncoding() *EncodingTest {
	byts, err := testingutils.TestSignedExecutionPayloadBid.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	root, err := testingutils.TestSignedExecutionPayloadBid.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"signed execution payload bid encoding",
		testdoc.SignedExecutionPayloadBidEncodingTestDoc,
		byts,
		root,
	)
}
