package aggregatorcommitteeconsensusdata

import (
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// AggregatorNoJustifications tests an invalid consensus data with no aggregator pre-consensus justifications
func AggregatorNoJustifications() *AggregatorCommitteeConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	return NewAggregatorCommitteeConsensusDataTest(
		"aggregator without justification",
		testdoc.AggregatorCommitteeConsensusDataTestAggregatorNoJustificationsDoc,
		*testingutils.TestAggregatorConsensusData(gloas.DataVersionGloas),
		0,
	)
}
