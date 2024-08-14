package runnerconstruction

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func ManyShares() tests.SpecTest {

	ks := testingutils.KeySetMapForValidators(10)
	shares := testingutils.ShareMapFromKeySetMap(ks)

	// No errors since one share must be valid for all runners
	expectedErrors := map[types.RunnerRole]string{
		types.RoleCommittee:                 "",
		types.RoleProposer:                  "must have one share",
		types.RoleAggregator:                "must have one share",
		types.RoleSyncCommitteeContribution: "must have one share",
		types.RoleValidatorRegistration:     "must have one share",
		types.RoleVoluntaryExit:             "must have one share",
	}

	return &RunnerConstructionSpecTest{
		Name:      "many shares",
		Shares:    shares,
		RoleError: expectedErrors,
	}
}
