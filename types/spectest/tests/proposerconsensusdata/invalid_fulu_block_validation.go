package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidFuluBlockValidation tests an invalid consensus data with fulu block
func InvalidFuluBlockValidation() *ProposerConsensusDataTest {
	version := spec.DataVersionFulu

	cd := &types.ProposerConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewProposerConsensusDataTest(
		"invalid fulu block",
		testdoc.ProposerConsensusDataTestInvalidFuluBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
