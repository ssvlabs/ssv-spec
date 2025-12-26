package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidDenebBlindedBlockValidation tests an invalid consensus data with deneb blinded block
func InvalidDenebBlindedBlockValidation() *ProposerConsensusDataTest {
	version := spec.DataVersionDeneb

	cd := &types.ProposerConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewProposerConsensusDataTest(
		"invalid deneb blinded block",
		testdoc.ProposerConsensusDataTestInvalidDenebBlindedBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
