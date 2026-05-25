package aggregatorcommitteeconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ConsensusDataEncoding tests encoding and decoding ProposerConsensusData for all duties
func ConsensusDataEncoding(name, documentation string, cd *types.AggregatorCommitteeConsensusData) *EncodingTest {

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

func Phase0AggregatorConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding(
		"phase0 aggregation encoding",
		testdoc.AggregatorCommitteeConsensusDataEncodingTestPhase0AggregatorDoc,
		testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0),
	)
}
func ElectraAggregatorConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding(
		"electra aggregation encoding",
		testdoc.AggregatorCommitteeConsensusDataEncodingTestElectraAggregatorDoc,
		testingutils.TestAggregatorConsensusData(spec.DataVersionElectra),
	)
}
func SyncCommitteeContributionConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding(
		"sync committee contribution encoding",
		testdoc.AggregatorCommitteeConsensusDataEncodingTestSyncCommitteeContributionDoc,
		testingutils.TestSyncCommitteeContributionConsensusData,
	)
}
