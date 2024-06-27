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

	return &ValidatorConsensusDataTest{
		Name:          "validator registration",
		ConsensusData: cd,
		ExpectedError: "validator registration has no consensus data",
	}
}
