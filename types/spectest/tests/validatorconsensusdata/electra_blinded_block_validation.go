package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ElectraBlindedBlockValidation tests a valid consensus data with electra blinded block
func ElectraBlindedBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid electra blinded block",
		testdoc.ValidatorConsensusDataTestElectraBlindedBlockDoc,
		*testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionElectra),
		"",
	)
}
