package valcheckattestations

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests valid data
func Valid() tests.SpecTest {
	return &valcheck.SpecTest{
		Name:       "attestation value check valid",
		Network:    types.PraterNetwork,
		RunnerRole: types.RoleCommittee,
		Input:      testingutils.TestBeaconVoteByts,
	}
}
