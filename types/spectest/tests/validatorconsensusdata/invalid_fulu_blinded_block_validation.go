package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidFuluBlindedBlockValidation tests an invalid consensus data with fulu blinded block
func InvalidFuluBlindedBlockValidation() *ValidatorConsensusDataTest {
	version := spec.DataVersionFulu

	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewValidatorConsensusDataTest(
		"invalid fulu blinded block",
		testdoc.ValidatorConsensusDataTestInvalidFuluBlindedBlockDoc,
		*cd,
		"could not unmarshal ssz (blinded err: incorrect size, regular err: incorrect size)",
	)
}
