package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ElectraBlockValidation tests a valid consensus data with electra block
func ElectraBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid electra block",
		"Test validation of valid consensus data with Electra block",
		*testingutils.TestProposerConsensusDataV(spec.DataVersionElectra),
		"",
	)
}
