package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ElectraBlockValidation tests a valid consensus data with electra block
func ElectraBlockValidation() *ValidatorConsensusDataTest {
	return &ValidatorConsensusDataTest{
		Name:          "valid electra block",
		ConsensusData: *testingutils.TestProposerConsensusDataV(spec.DataVersionElectra),
	}
}
