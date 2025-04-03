package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ElectraBlindedBlockValidation tests a valid consensus data with electra blinded block
func ElectraBlindedBlockValidation() *ValidatorConsensusDataTest {
	return &ValidatorConsensusDataTest{
		Name:          "valid electra blinded block",
		ConsensusData: *testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionElectra),
	}
}
