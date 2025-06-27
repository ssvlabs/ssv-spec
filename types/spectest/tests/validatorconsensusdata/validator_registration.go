package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidatorRegistration tests an invalid consensus data for validator registration (has no consensus data)
func ValidatorRegistration() *ValidatorConsensusDataTest {

	dataByts, err := testingutils.TestingValidatorRegistration.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	cd := types.ValidatorConsensusData{
		Duty:    testingutils.TestingValidatorRegistrationDuty,
		Version: spec.DataVersionCapella,
		DataSSZ: dataByts,
	}

	return NewValidatorConsensusDataTest(
		"validator registration",
		"Test validation error for validator registration consensus data which has no consensus data",
		cd,
		"validator registration has no consensus data",
	)
}
