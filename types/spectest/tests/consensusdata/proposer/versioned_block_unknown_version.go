package consensusdataproposer

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// VersionedBlockUnknownVersion tests a valid consensus data with unknown block
func VersionedBlockUnknownVersion() *ProposerSpecTest {
	unknownDataVersion := spec.DataVersion(100)
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerConsensusDataV(ks, spec.DataVersionCapella)
	cd.Version = unknownDataVersion

	cdSSZ, err := cd.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return &ProposerSpecTest{
		Name:          "consensus data versioned block unknown version",
		DataCd:        cdSSZ,
		ExpectedError: fmt.Sprintf("unknown block version %s", unknownDataVersion.String()),
	}
}
