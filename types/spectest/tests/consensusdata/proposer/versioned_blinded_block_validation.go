package consensusdataproposer

import "github.com/bloxapp/ssv-spec/types/testingutils"

// VersionedBlindedBlockValidation tests a valid consensus data with capella blinded block
func VersionedBlindedBlockValidation() *ProposerSpecTest {
	expectedCdRoot, err := testingutils.TestProposerBlindedBlockConsensusData.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	expectedBlkRoot, err := testingutils.TestingBlindedBeaconBlock.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &ProposerSpecTest{
		Name:            "consensus data versioned blinded block validation",
		Blinded:         true,
		DataCd:          testingutils.TestProposerBlindedBlockConsensusDataByts,
		DataBlk:         testingutils.TestingBlindedBeaconBlockBytes,
		ExpectedCdRoot:  expectedCdRoot,
		ExpectedBlkRoot: expectedBlkRoot,
	}
}
