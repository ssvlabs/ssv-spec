package runnerconstruction

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func OneShare() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	shares := testingutils.ShareMapFromKeySetMap(map[phase0.ValidatorIndex]*testingutils.TestKeySet{
		testingutils.TestingValidatorIndex: ks,
	})

	// No errors since one share must be valid for all runners
	expectedErrors := map[types.RunnerRole]string{
		types.RoleCommittee:                 "",
		types.RoleProposer:                  "",
		types.RoleAggregator:                "",
		types.RoleSyncCommitteeContribution: "",
		types.RoleValidatorRegistration:     "",
		types.RoleVoluntaryExit:             "",
	}

	return NewRunnerConstructionSpecTest(
		"one share",
		"Test that all runners can be constructed with one share",
		shares,
		expectedErrors,
	)
}
