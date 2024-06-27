package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// VoluntaryExit tests an invalid consensus data for voluntary exit (has no consensus data)
func VoluntaryExit() *ValidatorConsensusDataTest {

	dataByts, err := testingutils.TestingVoluntaryExit.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	cd := types.ValidatorConsensusData{
		Duty:    testingutils.TestingVoluntaryExitDuty,
		Version: spec.DataVersionCapella,
		DataSSZ: dataByts,
	}

	return &ValidatorConsensusDataTest{
		Name:          "voluntary exit",
		ConsensusData: cd,
		ExpectedError: "voluntary exit has no consensus data",
	}
}
