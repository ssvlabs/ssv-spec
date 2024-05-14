package valcheckcommittee

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests valid data
func Valid() tests.SpecTest {
	return &valcheck.SpecTest{
		Name:    "committee value check valid",
		Network: types.BeaconTestNetwork,
		Role:    types.RoleCommittee,
		Input:   testingutils.TestBeaconVoteByts,
	}
}
