package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type DutySpecTest struct {
	Name       string
	BeaconRole types.BeaconRole
	RunnerRole types.RunnerRole
}

func (test *DutySpecTest) TestName() string {
	return "duty " + test.Name
}

func (test *DutySpecTest) Run(t *testing.T) {
	result := types.MapDutyToRunnerRole(test.BeaconRole)
	assert.Equal(t, test.RunnerRole, result)
}
