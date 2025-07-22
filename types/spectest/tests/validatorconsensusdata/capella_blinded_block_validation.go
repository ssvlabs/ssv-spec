package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CapellaBlindedBlockValidation tests a valid consensus data with capella blinded block
func CapellaBlindedBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid capella blinded block",
		testdoc.ValidatorConsensusDataTestCapellaBlindedBlockDoc,
		*testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionCapella),
		"",
	)
}
