package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FuluBlockValidation tests a valid consensus data with fulu block
func FuluBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid fulu block",
		testdoc.ValidatorConsensusDataTestElectraBlockDoc,
		*testingutils.TestProposerConsensusDataV(spec.DataVersionFulu),
		0,
	)
}
