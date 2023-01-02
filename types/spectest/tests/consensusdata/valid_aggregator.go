package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidAggregator tests a valid aggregator consensus data
func ValidAggregator() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name: "valid aggregator",
		Obj: &types.ConsensusData{
			Duty:              testingutils.TestingAggregatorDuty,
			AggregateAndProof: testingutils.TestingAggregateAndProof,
		},
	}
}
