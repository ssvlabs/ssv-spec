package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CapellaBlockValidation tests a valid consensus data with capella block
func CapellaBlockValidation() *ValidatorConsensusDataTest {
	return &ValidatorConsensusDataTest{
		Name:          "valid capella block",
		ConsensusData: *testingutils.TestProposerConsensusDataV(spec.DataVersionCapella),
	}
}
