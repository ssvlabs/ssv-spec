package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
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
	return NewValidatorConsensusDataTest(
		"invalid capella block",
		testdoc.ValidatorConsensusDataTestInvalidCapellaBlockDoc,
		*cd,
		"could not unmarshal ssz (blinded err: incorrect size, regular err: incorrect size)",
	)
}
