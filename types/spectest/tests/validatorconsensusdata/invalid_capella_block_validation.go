package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidCapellaBlockValidation tests an invalid consensus data with capella block
func InvalidCapellaBlockValidation() *ValidatorConsensusDataTest {

	version := spec.DataVersionCapella

	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return &ValidatorConsensusDataTest{
		Name:          "invalid capella block",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: incorrect size",
	}
}
