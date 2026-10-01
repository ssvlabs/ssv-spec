package aggregatorcommitteeconsensusdata

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
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

func AggregatorConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding(
		"aggregation encoding",
		testdoc.AggregatorCommitteeConsensusDataEncodingTestAggregatorDoc,
		testingutils.TestAggregatorConsensusData(gloas.DataVersionGloas),
	)
}
func SyncCommitteeContributionConsensusDataEncoding() *EncodingTest {
	return ConsensusDataEncoding(
		"sync committee contribution encoding",
		testdoc.AggregatorCommitteeConsensusDataEncodingTestSyncCommitteeContributionDoc,
		testingutils.TestSyncCommitteeContributionConsensusData,
	)
}
