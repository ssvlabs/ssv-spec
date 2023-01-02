package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
