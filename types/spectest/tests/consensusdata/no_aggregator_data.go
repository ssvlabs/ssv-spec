package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
