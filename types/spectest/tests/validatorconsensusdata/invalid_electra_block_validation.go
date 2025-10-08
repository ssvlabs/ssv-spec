package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidElectraBlockValidation tests an invalid consensus data with electra block
func InvalidElectraBlockValidation() *ValidatorConsensusDataTest {

	version := spec.DataVersionElectra

	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewValidatorConsensusDataTest(
		"invalid electra block",
		testdoc.ValidatorConsensusDataTestInvalidElectraBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
