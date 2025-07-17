package runnerconstruction

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func ManyShares() tests.SpecTest {

	ks := testingutils.KeySetMapForValidators(10)
	shares := testingutils.ShareMapFromKeySetMap(ks)

	expectedErrors := map[types.RunnerRole]string{
		types.RoleCommittee:                 "", // No errors since it can handle multiple shares
		types.RoleProposer:                  "must have one share",
		types.RoleAggregator:                "must have one share",
		types.RoleSyncCommitteeContribution: "must have one share",
		types.RoleValidatorRegistration:     "must have one share",
		types.RoleVoluntaryExit:             "must have one share",
	}

	return NewRunnerConstructionSpecTest(
		"many shares",
		"Test that only committee runner can be constructed with multiple shares",
		shares,
		expectedErrors,
	)
}
