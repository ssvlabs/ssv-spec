package duty

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/assert"
)

type DutySpecTest struct {
	Name          string
	Type          string
	Documentation string
	BeaconRole    types.BeaconRole
	RunnerRole    types.RunnerRole
}

func (test *DutySpecTest) TestName() string {
	return "duty " + test.Name
}

func (test *DutySpecTest) Run(t *testing.T) {
	result := types.MapDutyToRunnerRole(test.BeaconRole)
	assert.Equal(t, test.RunnerRole, result)
}

func NewDutySpecTest(name, documentation string, beaconRole types.BeaconRole, runnerRole types.RunnerRole) *DutySpecTest {
	return &DutySpecTest{
		Name:          name,
		Type:          "Duty: validation of beacon role to runner role conversion",
		Documentation: documentation,
		BeaconRole:    beaconRole,
		RunnerRole:    runnerRole,
	}
}
