package consensusdataproposer

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// VersionedBlockValidation tests a valid consensus data with Deneb block
func VersionedBlockValidation() *ProposerSpecTest {
	expectedCdRoot, err := testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb).HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	expectedBlkRoot, err := testingutils.TestingBeaconBlockV(spec.DataVersionDeneb).Root()
	if err != nil {
		panic(err.Error())
	}

	return NewProposerSpecTest(
		"consensus data versioned block validation",
		testdoc.ProposerSpecTestVersionedBlockValidationDoc,
		false,
		testingutils.TestProposerConsensusDataBytsV(spec.DataVersionDeneb),
		testingutils.TestingBeaconBlockBytesV(spec.DataVersionDeneb),
		types.ExpectedBlkRoot(expectedBlkRoot),
		types.ExpectedCdRoot(expectedCdRoot),
		"",
	)
}
