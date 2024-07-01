package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DenebBlindedBlockValidation tests a valid consensus data with deneb blinded block
func DenebBlindedBlockValidation() *ValidatorConsensusDataTest {
	return &ValidatorConsensusDataTest{
		Name:          "valid deneb blinded block",
		ConsensusData: *testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb),
	}
}
