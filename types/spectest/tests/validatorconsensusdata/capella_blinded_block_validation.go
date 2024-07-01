package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CapellaBlindedBlockValidation tests a valid consensus data with capella blinded block
func CapellaBlindedBlockValidation() *ValidatorConsensusDataTest {
	ks := testingutils.Testing4SharesSet()

	return &ValidatorConsensusDataTest{
		Name:          "valid capella blinded block",
		ConsensusData: *testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionCapella),
	}
}
