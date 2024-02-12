package consensusdata

import "github.com/bloxapp/ssv-spec/types/testingutils"

// AggregatorNoJustifications tests an invalid consensus data with no aggregator pre-consensus justifications
func AggregatorNoJustifications() *ConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	return &ConsensusDataTest{
		Name:          "aggregator without justification",
		ConsensusData: *testingutils.TestAggregatorConsensusData,
	}
}
