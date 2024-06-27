package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidDenebBlockValidation tests an invalid consensus data with deneb block
func InvalidDenebBlockValidation() *ValidatorConsensusDataTest {
	version := spec.DataVersionDeneb
	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return &ValidatorConsensusDataTest{
		Name:          "invalid deneb block",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: incorrect size",
	}
}
