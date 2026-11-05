package proposerconsensusdata

import (
	"github.com/ssvlabs/ssv-spec/types/gloas"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidDuty tests an invalid consensus data with invalid duty
func InvalidDuty() *ProposerConsensusDataTest {

	cd := &types.ProposerConsensusData{
		Duty: types.ValidatorDuty{
			Type:   types.BeaconRole(100),
			PubKey: testingutils.TestingValidatorPubKey,
		},
		Version: gloas.DataVersionGloas,
		DataSSZ: testingutils.TestingBeaconBlockBytesV(gloas.DataVersionGloas),
	}

	return NewProposerConsensusDataTest(
		"invalid duty",
		testdoc.ProposerConsensusDataTestInvalidDutyDoc,
		*cd,
		types.UnknownDutyRoleDataErrorCode,
	)
}
