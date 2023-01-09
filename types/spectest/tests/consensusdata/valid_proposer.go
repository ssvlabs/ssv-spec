package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidProposer tests a valid block proposer
func ValidProposer() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name: "valid proposer",
		Obj: &types.ConsensusData{
			Duty:      testingutils.TestingProposerDuty,
			BlockData: testingutils.TestingBeaconBlock,
		},
	}
}
