package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongDutyTypeValidatorRegistration tests an invalid consensus data for validator registration (has no consensus data)
func WrongDutyTypeValidatorRegistration() *ProposerConsensusDataTest {

	dataByts, err := testingutils.TestingValidatorRegistration.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	cd := types.ProposerConsensusData{
		Duty:    testingutils.TestingValidatorRegistrationDuty,
		Version: spec.DataVersionCapella,
		DataSSZ: dataByts,
	}

	return NewProposerConsensusDataTest(
		"wrong duty type validator registration",
		testdoc.ProposerConsensusDataTestValidatorRegistrationDoc,
		cd,
		types.UnknownDutyRoleDataErrorCode,
	)
}

// WrongDutyTypeVoluntaryExit tests an invalid consensus data for voluntary exit (has no consensus data)
func WrongDutyTypeVoluntaryExit() *ProposerConsensusDataTest {

	dataByts, err := testingutils.TestingVoluntaryExit.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	cd := types.ProposerConsensusData{
		Duty:    testingutils.TestingVoluntaryExitDuty,
		Version: spec.DataVersionCapella,
		DataSSZ: dataByts,
	}

	return NewProposerConsensusDataTest(
		"wrong duty type voluntary exit",
		testdoc.ProposerConsensusDataTestVoluntaryExitDoc,
		cd,
		types.UnknownDutyRoleDataErrorCode,
	)
}
