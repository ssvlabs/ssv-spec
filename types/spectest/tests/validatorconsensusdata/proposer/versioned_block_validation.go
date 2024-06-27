package consensusdataproposer

import (
	"github.com/attestantio/go-eth2-client/spec"

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

	return &ProposerSpecTest{
		Name:            "consensus data versioned block validation",
		DataCd:          testingutils.TestProposerConsensusDataBytsV(spec.DataVersionDeneb),
		DataBlk:         testingutils.TestingBeaconBlockBytesV(spec.DataVersionDeneb),
		ExpectedCdRoot:  expectedCdRoot,
		ExpectedBlkRoot: expectedBlkRoot,
	}
}
