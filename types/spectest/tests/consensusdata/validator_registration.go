package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidatorRegistration tests an invalid consensus data for validator registration (has no consensus data)
func ValidatorRegistration() *ConsensusDataTest {

	dataByts, err := testingutils.TestingValidatorRegistration.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	cd := types.ConsensusData{
		Duty:    testingutils.TestingValidatorRegistrationDuty,
		Version: spec.DataVersionCapella,
		DataSSZ: dataByts,
	}

	return &ConsensusDataTest{
		Name:          "validator registration",
		ConsensusData: cd,
		ExpectedError: "validator registration has no consensus data",
	}
}
