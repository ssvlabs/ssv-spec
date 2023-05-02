package consensusdataproposer

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// VersionedBlindedBlockUnknownVersion tests a valid consensus data with unknown block version
func VersionedBlindedBlockUnknownVersion() *ProposerSpecTest {
	unknownDataVersion := spec.DataVersion(100)
	cd := &types.ConsensusData{
		Duty:    testingutils.TestingProposerDuty,
		Version: unknownDataVersion,
		DataSSZ: testingutils.TestProposerBlindedBlockConsensusDataByts,
	}

	dataCd, err := cd.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return &ProposerSpecTest{
		Name:            "consensus data versioned blinded block unknown version",
		DataCd:          dataCd,
		DataBlk:         nil,
		ExpectedCdRoot:  [32]byte{},
		ExpectedBlkRoot: [32]byte{},
		ExpectedError:   fmt.Sprintf("unknown block version %s", unknownDataVersion.String()),
	}
}
