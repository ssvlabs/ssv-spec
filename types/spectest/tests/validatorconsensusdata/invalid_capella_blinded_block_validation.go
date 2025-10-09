package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidCapellaBlindedBlockValidation tests an invalid consensus data with capella blinded block
func InvalidCapellaBlindedBlockValidation() *ValidatorConsensusDataTest {
	version := spec.DataVersionCapella

	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewValidatorConsensusDataTest(
		"invalid capella blinded block",
		testdoc.ValidatorConsensusDataTestInvalidCapellaBlindedBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
