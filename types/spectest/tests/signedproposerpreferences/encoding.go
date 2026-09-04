package signedproposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignedProposerPreferencesEncoding tests encoding and decoding a SignedProposerPreferences object (SIP #94 §5).
func SignedProposerPreferencesEncoding() *EncodingTest {

	s := testingutils.TestSignedProposerPreferences

	byts, err := s.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	root, err := s.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"signed proposer preferences encoding",
		testdoc.SignedProposerPreferencesEncodingTestDoc,
		byts,
		root,
	)
}
