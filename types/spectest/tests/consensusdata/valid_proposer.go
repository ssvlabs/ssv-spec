package consensusdata

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
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
