package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidCapellaBlockValidation tests an invalid consensus data with capella block
func InvalidCapellaBlockValidation() *ProposerConsensusDataTest {

	version := spec.DataVersionCapella

	cd := &types.ProposerConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewProposerConsensusDataTest(
		"invalid capella block",
		testdoc.ProposerConsensusDataTestInvalidCapellaBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
