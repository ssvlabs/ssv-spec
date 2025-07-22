package consensusdataproposer

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// VersionedBlindedBlockValidation tests a valid consensus data with Deneb blinded block
func VersionedBlindedBlockValidation() *ProposerSpecTest {
	expectedCdRoot, err := testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb).HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	expectedBlkRoot, err := testingutils.TestingBlindedBeaconBlockV(spec.DataVersionDeneb).Root()
	if err != nil {
		panic(err.Error())
	}

	return NewProposerSpecTest(
		"consensus data versioned blinded block validation",
		testdoc.ProposerSpecTestVersionedBlindedBlockValidationDoc,
		true,
		testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionDeneb),
		testingutils.TestingBlindedBeaconBlockBytesV(spec.DataVersionDeneb),
		types.ExpectedCdRoot(expectedCdRoot),
		types.ExpectedBlkRoot(expectedBlkRoot),
		"",
	)
}
