package validatorconsensusdata

import "github.com/ssvlabs/ssv-spec/types/testingutils"

// AggregatorValidation tests a valid consensus data with AggregateAndProof
func AggregatorValidation() *ValidatorConsensusDataTest {
	return &ValidatorConsensusDataTest{
		Name:          "aggregator valid",
		ConsensusData: *testingutils.TestAggregatorConsensusData,
	}
}
