package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidFuluBlockValidation tests an invalid consensus data with fulu block
func InvalidFuluBlockValidation() *ValidatorConsensusDataTest {
	version := spec.DataVersionFulu

	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewValidatorConsensusDataTest(
		"invalid fulu block",
		testdoc.ValidatorConsensusDataTestInvalidFuluBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
