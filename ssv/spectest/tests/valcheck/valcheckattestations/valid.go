package valcheckattestations

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Valid tests valid data
func Valid() *AttestationValCheckSpecTest {
	return &AttestationValCheckSpecTest{
		Name:    "attestation value check valid",
		Network: types.PraterNetwork,
		Input:   testingutils.TestAttesterConsensusDataByts,
	}
}
