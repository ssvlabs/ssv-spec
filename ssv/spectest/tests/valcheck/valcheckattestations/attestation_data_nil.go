package valcheckattestations

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
		ExpectedError: "attestation data nil",
	}
}
