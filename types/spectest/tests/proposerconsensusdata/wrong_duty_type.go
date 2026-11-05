package proposerconsensusdata

import (
	"github.com/ssvlabs/ssv-spec/types/gloas"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongDutyTypeVoluntaryExit tests an invalid consensus data for voluntary exit (has no consensus data)
func WrongDutyTypeVoluntaryExit() *ProposerConsensusDataTest {

	dataByts, err := testingutils.TestingVoluntaryExit.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	cd := types.ProposerConsensusData{
		Duty:    testingutils.TestingVoluntaryExitDuty,
		Version: gloas.DataVersionGloas,
		DataSSZ: dataByts,
	}

	return NewProposerConsensusDataTest(
		"wrong duty type voluntary exit",
		testdoc.ProposerConsensusDataTestVoluntaryExitDoc,
		cd,
		types.UnknownDutyRoleDataErrorCode,
	)
}
