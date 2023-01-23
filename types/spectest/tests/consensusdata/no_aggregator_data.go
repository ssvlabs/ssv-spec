package consensusdata

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// NoAggregatorData tests a nil aggregator data
func NoAggregatorData() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name: "no aggregator data",
		Obj: &types.ConsensusData{
			Duty: testingutils.TestingAggregatorDuty,
		},
		ExpectedErr: "aggregate and proof data is nil",
	}
}
