package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ConsensusDataEncoding tests encoding and decoding ProposerConsensusData for all duties
func ConsensusDataEncoding(name, documentation string, cd *types.ProposerConsensusData) *EncodingTest {

	byts, err := cd.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := cd.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		name,
		documentation,
		byts,
		root,
	)
}

func ProposerConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding(
		"proposer encoding",
		testdoc.ProposerConsensusDataEncodingTestProposerDoc,
		testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionCapella),
	)
}
func BlindedProposerConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding(
		"blinded proposer encoding",
		testdoc.ProposerConsensusDataEncodingTestBlindedProposerDoc,
		testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionCapella),
	)
}
