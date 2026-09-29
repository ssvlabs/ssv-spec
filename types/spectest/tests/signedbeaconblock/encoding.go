package signedbeaconblock

import (
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignedBeaconBlockEncoding tests encoding and decoding a SignedBeaconBlock object.
func SignedBeaconBlockEncoding() *EncodingTest {
	byts, err := testingutils.TestSignedBeaconBlock.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	root, err := testingutils.TestSignedBeaconBlock.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"signed beacon block encoding",
		testdoc.SignedBeaconBlockEncodingTestDoc,
		byts,
		root,
	)
}

// SignedBeaconBlockDevnet6Encoding tests encoding and decoding a real Gloas SignedBeaconBlock captured from lighthouse v8.2.0
// (glamsterdam-devnet-6): the codec must byte-round-trip a real CL's wire format (SIP #94 §4).
func SignedBeaconBlockDevnet6Encoding() *EncodingTest {
	decoded := &gloas.SignedBeaconBlock{}
	if err := decoded.UnmarshalSSZ(gloas.TestingDevnet6SignedBeaconBlockSSZ); err != nil {
		panic(err.Error())
	}
	root, err := decoded.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"signed beacon block devnet-6 encoding",
		testdoc.SignedBeaconBlockDevnet6EncodingTestDoc,
		gloas.TestingDevnet6SignedBeaconBlockSSZ,
		root,
	)
}
