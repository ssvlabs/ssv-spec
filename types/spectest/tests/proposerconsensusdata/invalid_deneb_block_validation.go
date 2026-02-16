package proposerconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidDenebBlockValidation tests an invalid consensus data with deneb block
func InvalidDenebBlockValidation() *ProposerConsensusDataTest {
	version := spec.DataVersionDeneb
	cd := &types.ProposerConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return NewProposerConsensusDataTest(
		"invalid deneb block",
		testdoc.ProposerConsensusDataTestInvalidDenebBlockDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
