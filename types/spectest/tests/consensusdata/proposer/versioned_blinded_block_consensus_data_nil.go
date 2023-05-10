package consensusdataproposer

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// VersionedBlindedBlockConsensusDataNil tests an invalid consensus data with bellatrix block
func VersionedBlindedBlockConsensusDataNil() *ProposerSpecTest {
	cd := &types.ConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
		Version: spec.DataVersionBellatrix,
		DataSSZ: nil,
	}

	cdSSZ, err := cd.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return &ProposerSpecTest{
		Name:            "consensus data versioned blinded block corrupted consensus data",
		DataCd:          cdSSZ,
		DataBlk:         nil,
		ExpectedCdRoot:  [32]byte{},
		ExpectedBlkRoot: [32]byte{},
		ExpectedError:   "could not unmarshal ssz: incorrect size",
	}
}
