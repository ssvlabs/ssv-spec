package aggregatorcommitteeconsensusdata

import (
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// AggregatorValidation tests a valid consensus data with AggregateAndProof
func AggregatorValidation() *AggregatorCommitteeConsensusDataTest {
	return NewAggregatorCommitteeConsensusDataTest(
		"aggregator valid",
		testdoc.AggregatorCommitteeConsensusDataTestAggregatorValidationDoc,
		*testingutils.TestAggregatorConsensusData(gloas.DataVersionGloas),
		0,
	)
}
