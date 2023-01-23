package consensusdata

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
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
