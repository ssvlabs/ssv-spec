package consensusdataproposer

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// VersionedBlockValidation tests a valid consensus data with bellatrix block
func VersionedBlockValidation() *ProposerSpecTest {
	expectedCdRoot, err := testingutils.TestProposerConsensusDataV(spec.DataVersionBellatrix).HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	expectedBlkRoot, err := testingutils.TestingBeaconBlockV(spec.DataVersionBellatrix).Root()
	if err != nil {
		panic(err.Error())
	}

	return &ProposerSpecTest{
		Name:            "consensus data versioned block validation",
		DataCd:          testingutils.TestProposerConsensusDataBytsV(spec.DataVersionBellatrix),
		DataBlk:         testingutils.TestingBeaconBlockBytesV(spec.DataVersionBellatrix),
		ExpectedCdRoot:  expectedCdRoot,
		ExpectedBlkRoot: expectedBlkRoot,
	}
}
