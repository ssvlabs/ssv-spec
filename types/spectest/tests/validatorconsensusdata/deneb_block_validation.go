package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DenebBlockValidation tests a valid consensus data with deneb block
func DenebBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid deneb block",
		testdoc.ValidatorConsensusDataTestDenebBlockDoc,
		*testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb),
		"",
	)
}
