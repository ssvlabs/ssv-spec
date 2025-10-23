package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidElectraBlindedBlockValidation tests an invalid consensus data with electra blinded block
func InvalidElectraBlindedBlockValidation() *ValidatorConsensusDataTest {

	version := spec.DataVersionElectra

	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewValidatorConsensusDataTest(
		"invalid electra blinded block",
		testdoc.ValidatorConsensusDataTestInvalidElectraBlindedBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
