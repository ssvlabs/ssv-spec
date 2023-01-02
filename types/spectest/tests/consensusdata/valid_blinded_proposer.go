package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidBlindedProposer tests a valid blinded block proposer
func ValidBlindedProposer() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name: "valid blinded proposer",
		Obj: &types.ConsensusData{
			Duty:             testingutils.TestingProposerDuty,
			BlindedBlockData: testingutils.TestingBlindedBeaconBlock,
		},
	}
}
