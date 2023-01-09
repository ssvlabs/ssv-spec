package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
