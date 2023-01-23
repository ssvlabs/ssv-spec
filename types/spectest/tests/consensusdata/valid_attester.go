package consensusdata

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// ValidAttester tests a valid attester consensus data
func ValidAttester() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name: "valid attester",
		Obj: &types.ConsensusData{
			Duty:            testingutils.TestingAttesterDuty,
			AttestationData: testingutils.TestingAttestationData,
		},
	}
}
