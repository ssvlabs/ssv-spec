package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CapellaBlindedBlockValidation tests a valid consensus data with capella blinded block
func CapellaBlindedBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid capella blinded block",
		"Test validation of valid consensus data with Capella blinded block",
		*testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionCapella),
		"",
	)
}
