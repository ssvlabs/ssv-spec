package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
)

// DutyNil tests a nil duty obj
func DutyNil() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name:        "duty nil",
		Obj:         &types.ConsensusData{},
		ExpectedErr: "duty is nil",
	}
}
