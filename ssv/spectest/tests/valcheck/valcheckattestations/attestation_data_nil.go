package valcheckattestations

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ConsensusDataNil tests consensus data != nil
func ConsensusDataNil() tests.SpecTest {
	consensusData := &types.ConsensusData{
		Duty:    testingutils.TestingAttesterDuty,
		DataSSZ: nil,
		Version: spec.DataVersionPhase0,
	}
	input, _ := consensusData.Encode()

	return &valcheck.SpecTest{
		Name:          "consensus data value check nil",
		Network:       types.PraterNetwork,
		BeaconRole:    types.BNRoleAttester,
		Input:         input,
		ExpectedError: "invalid value: could not unmarshal ssz: incorrect size",
	}
}
