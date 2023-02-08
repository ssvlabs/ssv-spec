package consensusdata

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

// DutyNil tests a nil duty obj
func DutyNil() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name:        "duty nil",
		Obj:         &types.ConsensusData{},
		ExpectedErr: "duty is nil",
	}
}
