package consensusdata

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
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
