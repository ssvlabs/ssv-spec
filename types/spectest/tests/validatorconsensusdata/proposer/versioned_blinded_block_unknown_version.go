package consensusdataproposer

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
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

	return &ProposerSpecTest{
		Name:          "consensus data versioned blinded block unknown version",
		DataCd:        dataCd,
		ExpectedError: fmt.Sprintf("unknown block version %s", unknownDataVersion.String()),
	}
}
