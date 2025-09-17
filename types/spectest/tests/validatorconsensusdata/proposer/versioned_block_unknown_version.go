package consensusdataproposer

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// VersionedBlockUnknownVersion tests a valid consensus data with unknown block
func VersionedBlockUnknownVersion() *ProposerSpecTest {
	unknownDataVersion := spec.DataVersion(100)
	cd := &types.ValidatorConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
		Version: unknownDataVersion,
		DataSSZ: testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionDeneb),
	}

	cdSSZ, err := cd.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return NewProposerSpecTest(
		"consensus data versioned block unknown version",
		testdoc.ProposerSpecTestVersionedBlockUnknownVersionDoc,
		false,
		cdSSZ,
		nil,
		[32]byte{},
		[32]byte{},
		fmt.Sprintf("unknown block version %d", unknownDataVersion),
	)
}
