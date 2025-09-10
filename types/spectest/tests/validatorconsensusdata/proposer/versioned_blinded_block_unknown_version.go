package consensusdataproposer

import (
    "github.com/attestantio/go-eth2-client/spec"

    "github.com/ssvlabs/ssv-spec/types"
    "github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
    "github.com/ssvlabs/ssv-spec/types/spectest/tests/errcodes"
    "github.com/ssvlabs/ssv-spec/types/testingutils"
)

// VersionedBlindedBlockUnknownVersion tests a valid consensus data with unknown block version
func VersionedBlindedBlockUnknownVersion() *ProposerSpecTest {
	unknownDataVersion := spec.DataVersion(100)
	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
		Version: unknownDataVersion,
		DataSSZ: testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionDeneb),
	}

	dataCd, err := cd.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return NewProposerSpecTest(
		"consensus data versioned blinded block unknown version",
		testdoc.ProposerSpecTestVersionedBlindedBlockUnknownVersionDoc,
		false,
		dataCd,
		nil,
		[32]byte{},
		[32]byte{},
		errcodes.ErrUnknownBlockVersion,
	)
}
