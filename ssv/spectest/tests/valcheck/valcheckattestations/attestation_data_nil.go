package valcheckattestations

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// AttestationDataNil tests attestation data != nil
func AttestationDataNil() *valcheck.SpecTest {
	cd := &types.ConsensusData{
		Duty:            testingutils.TestingAttesterDuty,
		AttestationData: nil,
	}
	source, _ := cd.MarshalSSZ()
	root, _ := cd.HashTreeRoot()
	input := &qbft.Data{
		Root:   root,
		Source: source,
	}

	return &valcheck.SpecTest{
		Name:          "attestation value check data nil",
		Network:       types.PraterNetwork,
		BeaconRole:    types.BNRoleAttester,
		Input:         input,
		ExpectedError: "failed decoding consensus data: incorrect size",
	}
}
