package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CapellaBlockValidation tests a valid consensus data with capella block
func CapellaBlockValidation() *ValidatorConsensusDataTest {
	ks := testingutils.Testing4SharesSet()

	return &ValidatorConsensusDataTest{
		Name:          "valid capella block",
		ConsensusData: *testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionCapella),
	}
}
