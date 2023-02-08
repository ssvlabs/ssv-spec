package consensusdata

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// NilBlock tests a nil block proposer
func NilBlock() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name: "nil block",
		Obj: &types.ConsensusData{
			Duty: testingutils.TestingProposerDuty,
		},
		ExpectedErr: "block data is nil",
	}
}
