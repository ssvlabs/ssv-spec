package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DenebBlockValidation tests a valid consensus data with deneb block
func DenebBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid deneb block",
		"Test validation of valid consensus data with Deneb block",
		*testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb),
		"",
	)
}
