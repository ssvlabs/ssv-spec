package consensusdata

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// NoAttestationData tests an invalid attester consensus data
func NoAttestationData() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name: "no attestation data",
		Obj: &types.ConsensusData{
			Duty: testingutils.TestingAttesterDuty,
		},
		ExpectedErr: "attestation data is nil",
	}
}
