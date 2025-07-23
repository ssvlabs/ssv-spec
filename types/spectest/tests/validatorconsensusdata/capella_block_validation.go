package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CapellaBlockValidation tests a valid consensus data with capella block
func CapellaBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid capella block",
		testdoc.ValidatorConsensusDataTestCapellaBlockDoc,
		*testingutils.TestProposerConsensusDataV(spec.DataVersionCapella),
		"",
	)
}
