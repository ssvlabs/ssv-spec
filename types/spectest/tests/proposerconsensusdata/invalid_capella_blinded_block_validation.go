package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidCapellaBlindedBlockValidation tests an invalid consensus data with capella blinded block
func InvalidCapellaBlindedBlockValidation() *ProposerConsensusDataTest {
	version := spec.DataVersionCapella

	cd := &types.ProposerConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewProposerConsensusDataTest(
		"invalid capella blinded block",
		testdoc.ProposerConsensusDataTestInvalidCapellaBlindedBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
