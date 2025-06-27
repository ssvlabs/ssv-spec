package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ElectraBlindedBlockValidation tests a valid consensus data with electra blinded block
func ElectraBlindedBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid electra blinded block",
		"Test validation of valid consensus data with Electra blinded block",
		*testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionElectra),
		"",
	)
}
