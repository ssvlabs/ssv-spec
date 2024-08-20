package runnerconstruction

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
)

func NoShares() tests.SpecTest {

	shares := map[phase0.ValidatorIndex]*types.Share{}

	expectedErrors := map[types.RunnerRole]string{
		types.RoleCommittee:                 "no shares",
		types.RoleProposer:                  "must have one share",
		types.RoleAggregator:                "must have one share",
		types.RoleSyncCommitteeContribution: "must have one share",
		types.RoleValidatorRegistration:     "must have one share",
		types.RoleVoluntaryExit:             "must have one share",
	}

	return &RunnerConstructionSpecTest{
		Name:      "no shares",
		Shares:    shares,
		RoleError: expectedErrors,
	}
}
