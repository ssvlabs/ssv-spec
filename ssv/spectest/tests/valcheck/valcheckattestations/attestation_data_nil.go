package valcheckattestations

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests/valcheck"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// AttestationDataNil tests attestation data != nil
func AttestationDataNil() *valcheck.SpecTest {
	consensusData := &types.ConsensusData{
		Duty:            testingutils.TestingAttesterDuty,
		AttestationData: nil,
	}
	input, _ := consensusData.Encode()

	return &valcheck.SpecTest{
		Name:          "attestation value check data nil",
		Network:       types.PraterNetwork,
		BeaconRole:    types.BNRoleAttester,
		Input:         input,
		ExpectedError: "invalid value: attestation data is nil",
	}
}
