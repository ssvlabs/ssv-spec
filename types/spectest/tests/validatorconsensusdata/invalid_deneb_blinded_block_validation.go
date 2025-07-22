package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidDenebBlindedBlockValidation tests an invalid consensus data with deneb blinded block
func InvalidDenebBlindedBlockValidation() *ValidatorConsensusDataTest {
	version := spec.DataVersionDeneb

	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewValidatorConsensusDataTest(
		"invalid deneb blinded block",
		testdoc.ValidatorConsensusDataTestInvalidDenebBlindedBlockDoc,
		*cd,
		"could not unmarshal ssz: incorrect size",
	)
}
