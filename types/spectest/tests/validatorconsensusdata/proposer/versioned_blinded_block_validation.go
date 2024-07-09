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

	return &ProposerSpecTest{
		Name:            "consensus data versioned blinded block validation",
		Blinded:         true,
		DataCd:          testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionDeneb),
		DataBlk:         testingutils.TestingBlindedBeaconBlockBytesV(spec.DataVersionDeneb),
		ExpectedCdRoot:  expectedCdRoot,
		ExpectedBlkRoot: expectedBlkRoot,
	}
}
