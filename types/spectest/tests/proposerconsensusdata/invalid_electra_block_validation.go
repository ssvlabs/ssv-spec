package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidElectraBlockValidation tests an invalid consensus data with electra block
func InvalidElectraBlockValidation() *ProposerConsensusDataTest {

	version := spec.DataVersionElectra

	cd := &types.ProposerConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewProposerConsensusDataTest(
		"invalid electra block",
		testdoc.ProposerConsensusDataTestInvalidElectraBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
