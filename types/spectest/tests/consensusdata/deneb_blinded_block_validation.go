package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DenebBlindedBlockValidation tests a valid consensus data with deneb blinded block
func DenebBlindedBlockValidation() *ValidatorConsensusDataTest {
	ks := testingutils.Testing4SharesSet()

	return &ValidatorConsensusDataTest{
		Name:          "valid deneb blinded block",
		ConsensusData: *testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb),
	}
}
