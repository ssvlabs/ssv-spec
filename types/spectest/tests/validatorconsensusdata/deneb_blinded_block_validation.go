package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DenebBlindedBlockValidation tests a valid consensus data with deneb blinded block
func DenebBlindedBlockValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"valid deneb blinded block",
		*testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb),
		"",
	)
}
