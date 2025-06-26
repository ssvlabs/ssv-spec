package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CapellaBlockValidation tests a valid consensus data with capella block
func CapellaBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid capella block",
		"Test validation of valid consensus data with Capella block",
		*testingutils.TestProposerConsensusDataV(spec.DataVersionCapella),
		"",
	)
}
