package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
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
	return NewValidatorConsensusDataTest(
		"invalid deneb block",
		testdoc.ValidatorConsensusDataTestInvalidDenebBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
