package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
