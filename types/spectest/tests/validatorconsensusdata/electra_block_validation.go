package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ElectraBlockValidation tests a valid consensus data with electra block
func ElectraBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid electra block",
		testdoc.ValidatorConsensusDataTestElectraBlockDoc,
		*testingutils.TestProposerConsensusDataV(spec.DataVersionElectra),
		0,
	)
}
