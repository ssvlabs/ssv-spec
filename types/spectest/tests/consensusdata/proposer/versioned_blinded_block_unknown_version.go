package consensusdataproposer

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// VersionedBlindedBlockUnknownVersion tests a valid consensus data with unknown block
func VersionedBlindedBlockUnknownVersion() *ProposerSpecTest {
	cd := &types.ConsensusData{
		Duty:    testingutils.TestingProposerDuty,
		Version: spec.DataVersionBellatrix,
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
		ExpectedError:   fmt.Sprintf("unknown block version %s", spec.DataVersionBellatrix),
	}
}
