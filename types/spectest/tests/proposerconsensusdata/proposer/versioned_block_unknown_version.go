package consensusdataproposer

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// VersionedBlockUnknownVersion tests a valid consensus data with unknown block
func VersionedBlockUnknownVersion() *ProposerSpecTest {
	unknownDataVersion := spec.DataVersion(100)
	cd := &types.ProposerConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(gloas.DataVersionGloas),
		Version: unknownDataVersion,
		DataSSZ: testingutils.TestingBeaconBlockBytesV(gloas.DataVersionGloas),
	}

	cdSSZ, err := cd.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return NewProposerSpecTest(
		"consensus data versioned block unknown version",
		testdoc.ProposerSpecTestVersionedBlockUnknownVersionDoc,
		cdSSZ,
		nil,
		[32]byte{},
		[32]byte{},
		types.UnknownBlockVersionErrorCode,
	)
}
