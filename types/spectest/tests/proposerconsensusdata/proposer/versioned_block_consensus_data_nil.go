package consensusdataproposer

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// VersionedBlockConsensusDataNil tests an invalid consensus data with nil Gloas proposal data
func VersionedBlockConsensusDataNil() *ProposerSpecTest {
	cd := &types.ProposerConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(gloas.DataVersionGloas),
		Version: gloas.DataVersionGloas,
		DataSSZ: nil,
	}

	cdSSZ, err := cd.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return NewProposerSpecTest(
		"consensus data versioned block corrupted consensus data",
		testdoc.ProposerSpecTestVersionedBlockConsensusDataNilDoc,
		cdSSZ,
		nil,
		[32]byte{},
		[32]byte{},
		types.UnmarshalSSZErrorCode,
	)
}
