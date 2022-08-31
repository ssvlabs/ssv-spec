package valcheckattestations

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// AttestationDataNil tests attestation data != nil
func AttestationDataNil() *AttestationValCheckSpecTest {
	consensusData := &types.ConsensusData{
		Duty:            testingutils.TestingAttesterDuty,
		AttestationData: nil,
	}
	input, _ := consensusData.Encode()

	return &AttestationValCheckSpecTest{
		Name:          "attestation value check data nil",
		Network:       types.PraterNetwork,
		Input:         input,
		ExpectedError: "attestation data nil",
	}
}
