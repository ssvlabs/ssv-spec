package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ConsensusDataEncoding tests encoding and decoding ValidatorConsensusData for all duties
func ConsensusDataEncoding(name, documentation string, cd *types.ValidatorConsensusData) *EncodingTest {

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
		testdoc.ValidatorConsensusDataEncodingTestProposerDoc,
		testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionCapella),
	)
}
func BlindedProposerConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding(
		"blinded proposer encoding",
		testdoc.ValidatorConsensusDataEncodingTestBlindedProposerDoc,
		testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionCapella),
	)
}
func Phase0AggregatorConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding(
		"phase0 aggregation encoding",
		testdoc.ValidatorConsensusDataEncodingTestPhase0AggregatorDoc,
		testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0),
	)
}
func ElectraAggregatorConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding(
		"electra aggregation encoding",
		testdoc.ValidatorConsensusDataEncodingTestElectraAggregatorDoc,
		testingutils.TestAggregatorConsensusData(spec.DataVersionElectra),
	)
}
func SyncCommitteeContributionConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding(
		"sync committee contribution encoding",
		testdoc.ValidatorConsensusDataEncodingTestSyncCommitteeContributionDoc,
		testingutils.TestSyncCommitteeContributionConsensusData,
	)
}
