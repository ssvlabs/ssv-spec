package consensusdataproposer

import "github.com/bloxapp/ssv-spec/types/testingutils"

// VersionedBlockValidation tests a valid consensus data with bellatrix block
func VersionedBlockValidation() *ProposerSpecTest {
	expectedCdRoot, err := testingutils.TestProposerConsensusData.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	expectedBlkRoot, err := testingutils.TestingBeaconBlock.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &ProposerSpecTest{
		Name:            "consensus data versioned block validation",
		DataCd:          testingutils.TestProposerConsensusDataByts,
		DataBlk:         testingutils.TestingBeaconBlockBytes,
		ExpectedCdRoot:  expectedCdRoot,
		ExpectedBlkRoot: expectedBlkRoot,
	}
}
