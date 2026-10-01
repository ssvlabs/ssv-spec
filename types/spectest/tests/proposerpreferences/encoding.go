package proposerpreferences

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ProposerPreferencesEncoding tests encoding and decoding a ProposerPreferences object (SIP #94 §5).
func ProposerPreferencesEncoding() *EncodingTest {

	p := testingutils.TestProposerPreferences

	byts, err := p.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	root, err := p.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"proposer preferences encoding",
		testdoc.ProposerPreferencesEncodingTestDoc,
		byts,
		root,
	)
}
