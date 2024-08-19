package consensusdataproposer

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// VersionedBlindedBlockConsensusDataNil tests an invalid consensus data with deneb block
func VersionedBlindedBlockConsensusDataNil() *ProposerSpecTest {
	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
		Version: spec.DataVersionDeneb,
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
