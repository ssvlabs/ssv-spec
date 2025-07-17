package consensusdataproposer

import (
	"github.com/attestantio/go-eth2-client/spec"

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
		"Test validation of valid consensus data with versioned Deneb blinded block",
		true,
		testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionDeneb),
		testingutils.TestingBlindedBeaconBlockBytesV(spec.DataVersionDeneb),
		expectedBlkRoot,
		expectedCdRoot,
		"",
	)
}
