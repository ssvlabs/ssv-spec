package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
)

type UnknownDuty struct {
}

func (ud *UnknownDuty) DutySlot() phase0.Slot {
	return 0
}

func (ud *UnknownDuty) RunnerRole() types.RunnerRole {
	return types.RoleCommittee
}
