package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidFuluBlindedBlockValidation tests an invalid consensus data with fulu blinded block
func InvalidFuluBlindedBlockValidation() *ProposerConsensusDataTest {
	version := spec.DataVersionFulu

	cd := &types.ProposerConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewProposerConsensusDataTest(
		"invalid fulu blinded block",
		testdoc.ProposerConsensusDataTestInvalidFuluBlindedBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
