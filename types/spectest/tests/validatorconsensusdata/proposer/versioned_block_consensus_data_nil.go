package consensusdataproposer

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// VersionedBlockConsensusDataNil tests an invalid consensus data with Deneb block
func VersionedBlockConsensusDataNil() *ProposerSpecTest {
	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
		Version: spec.DataVersionDeneb,
		DataSSZ: nil,
	}

	cdSSZ, err := cd.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return NewProposerSpecTest(
		"consensus data versioned block corrupted consensus data",
		testdoc.ProposerSpecTestVersionedBlockConsensusDataNilDoc,
		false,
		cdSSZ,
		nil,
		[32]byte{},
		[32]byte{},
		"could not unmarshal ssz (blinded err: incorrect size, regular err: incorrect size)",
	)
}
