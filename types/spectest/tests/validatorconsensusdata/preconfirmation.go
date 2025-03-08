package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Preconfirmation tests an invalid consensus data for preconfirmation (has no consensus data)
func Preconfirmation() *ValidatorConsensusDataTest {

	dataByts, err := testingutils.TestingPreconfRequest.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	cd := types.ValidatorConsensusData{
		Duty:    testingutils.TestingPreconfirmationDuty,
		Version: spec.DataVersionCapella,
		DataSSZ: dataByts,
	}

	return &ValidatorConsensusDataTest{
		Name:          "preconfirmation",
		ConsensusData: cd,
		ExpectedError: "preconfirmation has no consensus data",
	}
}
