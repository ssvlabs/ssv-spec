package validatorconsensusdata

import "github.com/ssvlabs/ssv-spec/types/testingutils"

// AggregatorNoJustifications tests an invalid consensus data with no aggregator pre-consensus justifications
func AggregatorNoJustifications() *ValidatorConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	return &ValidatorConsensusDataTest{
		Name:          "aggregator without justification",
		ConsensusData: *testingutils.TestAggregatorConsensusData,
	}
}
