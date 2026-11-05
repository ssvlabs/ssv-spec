package consensusdataproposer

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// VersionedBlockValidation tests a valid consensus data with a Gloas block
func VersionedBlockValidation() *ProposerSpecTest {
	version := gloas.DataVersionGloas

	expectedCdRoot, err := testingutils.TestProposerConsensusDataV(version).HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	block := gloas.TestingBeaconBlock(testingutils.TestingDutySlotV(version))
	expectedBlkRoot, err := block.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}
	blockSSZ, err := block.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return NewProposerSpecTest(
		"consensus data versioned block validation",
		testdoc.ProposerSpecTestVersionedBlockValidationDoc,
		testingutils.TestProposerConsensusDataBytsV(version),
		blockSSZ,
		types.ExpectedBlkRoot(expectedBlkRoot),
		types.ExpectedCdRoot(expectedCdRoot),
		0,
	)
}
