package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidDuty tests an invalid consensus data with invalid duty
func InvalidDuty() *ValidatorConsensusDataTest {

	cd := &types.ValidatorConsensusData{
		Duty: types.ValidatorDuty{
			Type:   types.BeaconRole(100),
			PubKey: testingutils.TestingValidatorPubKey,
		},
		Version: spec.DataVersionCapella,
		DataSSZ: testingutils.TestingAttestationDataBytes(spec.DataVersionCapella),
	}

	return NewValidatorConsensusDataTest(
		"invalid duty",
		testdoc.ValidatorConsensusDataTestInvalidDutyDoc,
		*cd,
		"unknown duty role",
	)
}
