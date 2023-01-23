package valcheckattestations

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests/valcheck"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Valid tests valid data
func Valid() *valcheck.SpecTest {
	return &valcheck.SpecTest{
		Name:       "attestation value check valid",
		Network:    types.PraterNetwork,
		BeaconRole: types.BNRoleAttester,
		Input:      testingutils.TestAttesterConsensusDataByts,
	}
}
