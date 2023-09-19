package consensusdataproposer

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// VersionedBlindedBlockValidation tests a valid consensus data with bellatrix blinded block
func VersionedBlindedBlockValidation() *ProposerSpecTest {
	expectedCdRoot, err := testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionBellatrix).HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	expectedBlkRoot, err := testingutils.TestingBlindedBeaconBlockV(spec.DataVersionBellatrix).Root()
	if err != nil {
		panic(err.Error())
	}

	return &ProposerSpecTest{
		Name:            "consensus data versioned blinded block validation",
		Blinded:         true,
		DataCd:          testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionBellatrix),
		DataBlk:         testingutils.TestingBlindedBeaconBlockBytesV(spec.DataVersionBellatrix),
		ExpectedCdRoot:  expectedCdRoot,
		ExpectedBlkRoot: expectedBlkRoot,
	}
}
