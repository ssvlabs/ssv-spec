package validator

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func InvalidDuty() tests.SpecTest {

	duty := testingutils.TestingAttesterDuty
	duty.Type = types.BeaconRole(100)

	return &ValidatorTest{
		Name:                   "invalid duty",
		Duties:                 []*types.Duty{&duty},
		Messages:               []*types.SSVMessage{},
		OutputMessages:         []*types.SSVMessage{},
		BeaconBroadcastedRoots: []string{},
		ExpectedError:          "duty type UNDEFINED not supported",
	}
}
